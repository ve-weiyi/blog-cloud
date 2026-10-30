package notification

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
)

// errContentGone 表示通知引用的内容已删除：这条投递不该再发，重试也不会成功。
var errContentGone = errors.New("notification: content gone")

const (
	// deliveryClaimBatch 单次领取的投递记录上限
	deliveryClaimBatch = 50
	// deliveryMaxAttempts 超过后置为失败终态，由管理端决定是否重发
	deliveryMaxAttempts = 5
	// deliveryBaseBackoff 指数退避的基数
	deliveryBaseBackoff = time.Minute
	// deliveryLease 领取后占用的租约：这段时间内别的实例不会再捞到同一条
	deliveryLease = 2 * time.Minute
)

// EmailSender 是通知域对外部邮件通道的唯一要求。
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, content string) error
}

// DeliveryWorker 把待发送的投递记录发出去。
//
// 投递表本身就是队列：定时领取一批待发送记录，发送后回写状态。
// 没有 MQ 唤醒也能工作，重试同样由这个循环驱动，不依赖消息重投。
type DeliveryWorker struct {
	svc      *Service
	sender   EmailSender
	interval time.Duration

	cancel context.CancelFunc
	done   chan struct{}
}

func NewDeliveryWorker(svc *Service, sender EmailSender, interval time.Duration) *DeliveryWorker {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &DeliveryWorker{svc: svc, sender: sender, interval: interval}
}

// Start 启动后台循环；重复调用无副作用。
func (w *DeliveryWorker) Start() {
	if w.cancel != nil {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	w.done = make(chan struct{})

	go func() {
		defer close(w.done)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.drain(ctx)
			}
		}
	}()
}

// Stop 停止循环并等待在途批次结束；nil 安全。
func (w *DeliveryWorker) Stop(ctx context.Context) {
	if w == nil || w.cancel == nil {
		return
	}
	w.cancel()
	select {
	case <-w.done:
	case <-ctx.Done():
	}
}

// drain 领取并处理一批；每条记录独立成败，不因一条失败中断整批。
func (w *DeliveryWorker) drain(ctx context.Context) {
	rows, err := w.svc.claimDeliveries(ctx, deliveryClaimBatch, deliveryLease)
	if err != nil {
		// 这批记录的租约到期后自动回队，不会丢；但必须留痕，否则积压时无从查起
		logx.WithContext(ctx).Errorf("领取投递记录失败: %v", err)
		return
	}
	if len(rows) == 0 {
		return
	}

	notifs, events, err := w.svc.loadDeliveryTargets(ctx, rows)
	if err != nil {
		logx.WithContext(ctx).Errorf("预取投递对象失败: %v", err)
		return
	}

	for _, row := range rows {
		notif := notifs[row.NotificationId]
		if notif == nil {
			// 条目已被撤回或删除，这条投递不该再发
			w.finishCancelled(ctx, row, "条目已删除")
			continue
		}
		if err := w.process(ctx, row, notif, events[notif.EventId]); err != nil {
			// 基础设施错误：保持待发送，下次循环再试
			w.retryLater(ctx, row, err)
		}
	}
}

// loadDeliveryTargets 批量取回本批投递引用的条目与事件。
//
// 逐条查会让往返次数随积压线性放大：一批 50 条是 100 次查询，批量只需 2 次。
func (s *Service) loadDeliveryTargets(ctx context.Context, rows []*model.TNotificationDelivery) (map[int64]*model.TNotification, map[int64]*model.TNotificationEvent, error) {
	ids := make([]int64, 0, len(rows))
	seen := make(map[int64]struct{}, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.NotificationId]; ok {
			continue
		}
		seen[row.NotificationId] = struct{}{}
		ids = append(ids, row.NotificationId)
	}

	notifs, err := s.loadEntries(ctx, ids)
	if err != nil {
		return nil, nil, err
	}
	events, err := s.loadEvents(ctx, notifs)
	if err != nil {
		return nil, nil, err
	}
	return notifs, events, nil
}

// process 发送一条投递。返回错误表示值得重试的基础设施问题。
func (w *DeliveryWorker) process(ctx context.Context, d *model.TNotificationDelivery, notif *model.TNotification, ev *model.TNotificationEvent) error {
	if ev == nil {
		return w.finish(ctx, d, DeliveryCancelled, "事件已删除")
	}
	if undeliverable(ev.Status) {
		return w.finish(ctx, d, DeliveryCancelled, "事件已终止")
	}

	scene, ok := notify.Lookup(notif.Scene)
	if !ok {
		return w.finish(ctx, d, DeliveryCancelled, "未知场景")
	}

	// 发送前复查偏好：记录可能已经排队几分钟，用户在这期间退订
	allowed, err := w.channelAllowed(ctx, notif.RecipientId, scene, notify.Channel(d.Channel))
	if err != nil {
		return err
	}
	if !allowed {
		return w.finish(ctx, d, DeliveryCancelled, "用户已退订")
	}

	title, body, err := w.render(ctx, notif, ev)
	if errors.Is(err, errContentGone) {
		return w.finish(ctx, d, DeliveryCancelled, "内容已删除")
	}
	if err != nil {
		return err
	}

	switch notify.Channel(d.Channel) {
	case notify.ChannelEmail:
		link := w.unsubscribeLink(notif.RecipientId, scene, notify.Channel(d.Channel))
		if err := w.sender.SendEmail(ctx, d.Recipient, title, emailShell(title, body, link)); err != nil {
			return err
		}
		return w.finish(ctx, d, DeliverySent, "")
	default:
		// 短信与推送尚未接入：置为取消，避免占着队列反复领取
		return w.finish(ctx, d, DeliveryCancelled, "渠道未接入")
	}
}

// render 渲染正文：公告读内容表，其余按场景与载荷渲染。
func (w *DeliveryWorker) render(ctx context.Context, notif *model.TNotification, ev *model.TNotificationEvent) (string, string, error) {
	payload := decodePayload(ev.Payload)

	if notif.Scene == notify.SceneAnnouncement {
		row, err := w.svc.announcements.FindById(ctx, payload.AnnouncementID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 公告已删：撤回路径本应同时回收条目，走到这里说明清理漏了，也不该再发
			return "", "", errContentGone
		}
		if err != nil {
			return "", "", err
		}
		return row.Title, deref(row.Content), nil
	}

	content, _ := notify.Render(notif.Scene, payload, 1, "")
	return content.Title, content.Body, nil
}

// channelAllowed 复查该渠道是否仍在偏好允许范围内。
//
// 强制分组不看偏好：偏好表本不该有它们的行，但脏数据不能拦下公告与账号安全。
func (w *DeliveryWorker) channelAllowed(ctx context.Context, userID string, scene notify.Scene, ch notify.Channel) (bool, error) {
	if scene.Mandatory {
		return true, nil
	}

	rows, err := w.svc.prefs.FindAll(ctx,
		"user_id = ? AND pref_group = ? AND channel IN ?",
		userID, string(scene.Group), []string{string(notify.ChannelInbox), string(ch)})
	if err != nil {
		return false, fmt.Errorf("notification: check preference: %w", err)
	}
	for _, row := range rows {
		if row.Enabled == 0 {
			return false, nil
		}
	}
	return true, nil
}

// finish 把记录置为终态。
func (w *DeliveryWorker) finish(ctx context.Context, d *model.TNotificationDelivery, status, cause string) error {
	fields := map[string]interface{}{
		"status":        status,
		"error_msg":     cause,
		"next_retry_at": nil,
	}
	if status == DeliverySent {
		fields["sent_at"] = time.Now()
	}
	_, err := w.svc.deliveries.UpdateFields(ctx, fields, "id = ?", d.Id)
	return err
}

// finishCancelled 把投递置为取消终态。
//
// 回写失败只记日志：记录仍是待发送，下一次领取会再走一遍同样的判断。
func (w *DeliveryWorker) finishCancelled(ctx context.Context, d *model.TNotificationDelivery, cause string) {
	if err := w.finish(ctx, d, DeliveryCancelled, cause); err != nil {
		logx.WithContext(ctx).Errorf("取消投递 %d 失败: %v", d.Id, err)
	}
}

// retryLater 记录失败原因并按指数退避安排下一次；超过上限置为失败终态。
func (w *DeliveryWorker) retryLater(ctx context.Context, d *model.TNotificationDelivery, cause error) {
	if d.Attempts >= deliveryMaxAttempts {
		if err := w.finish(ctx, d, DeliveryFailed, cause.Error()); err != nil {
			logx.WithContext(ctx).Errorf("投递 %d 置为失败终态失败: %v", d.Id, err)
		}
		return
	}

	attempt := d.Attempts
	if attempt < 1 {
		attempt = 1
	}
	backoff := deliveryBaseBackoff << (attempt - 1)
	if _, err := w.svc.deliveries.UpdateFields(ctx, map[string]interface{}{
		"status":        DeliveryPending,
		"error_msg":     cause.Error(),
		"next_retry_at": time.Now().Add(backoff),
	}, "id = ?", d.Id); err != nil {
		logx.WithContext(ctx).Errorf("安排投递 %d 重试失败: %v（失败原因: %v）", d.Id, err, cause)
	}
}

// claimDeliveries 领取一批待发送记录。
//
// 领取时写一个租约：这段时间内别的实例不会再捞到同一条，进程中途退出后租约到期自动重回队列。
func (s *Service) claimDeliveries(ctx context.Context, limit int, lease time.Duration) ([]*model.TNotificationDelivery, error) {
	var claimed []*model.TNotificationDelivery

	now := time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ids []int64
		if err := tx.Table(s.deliveries.TableName()).
			Select("id").
			Where("status = ? AND (next_retry_at IS NULL OR next_retry_at <= ?)", DeliveryPending, now).
			Order("id ASC").
			Limit(limit).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Pluck("id", &ids).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}

		if _, err := s.deliveries.WithTx(tx).UpdateFields(ctx, map[string]interface{}{
			"next_retry_at": now.Add(lease),
			"attempts":      gorm.Expr("attempts + 1"),
		}, "id IN ?", ids); err != nil {
			return err
		}

		rows, err := s.deliveries.WithTx(tx).FindAll(ctx, "id IN ?", ids)
		if err != nil {
			return err
		}
		claimed = rows
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("notification: claim deliveries: %w", err)
	}
	return claimed, nil
}

// undeliverable 判断事件是否已经不允许投递。
//
// 只认 terminated：done 是扇出正常完成，它的投递本来就该发；
// failed 意味着场景未知、根本没写条目，不会产生投递。不要复用 IsTerminal——
// 它服务的是「事件能否回收」，含 done，用在这里会把正常投递全部取消。
func undeliverable(status string) bool {
	return status == StatusTerminated
}

// unsubscribeTokenTTL 退订令牌有效期。够长（邮件可能很久才被点开）又不过期到失效链接永远可用。
const unsubscribeTokenTTL = 30 * 24 * time.Hour

// unsubscribeLink 为可退订分组签发退订链接。
//
// 强制分组不签发：它们的邮件里根本不该出现退订入口。密钥或站点地址缺失时返回空，
// 由调用方决定不渲染——拼一个点不开的链接比没有链接更糟。
func (w *DeliveryWorker) unsubscribeLink(userID string, scene notify.Scene, ch notify.Channel) string {
	// 密钥是否可用交给 SignUnsubscribe 裁决，这里只判断链接前缀
	if scene.Mandatory || w.svc.siteURL == "" {
		return ""
	}

	token, err := notify.SignUnsubscribe(w.svc.unsubscribeSecret, userID, scene.Group, ch,
		time.Now().Add(unsubscribeTokenTTL))
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s/unsubscribe?token=%s",
		strings.TrimRight(w.svc.siteURL, "/"), url.QueryEscape(token))
}

// emailShell 给正文套一层最简外壳，保持站内、邮件共用同一套文本。
// link 为空时不渲染退订页脚。
func emailShell(title, body, link string) string {
	footer := ""
	if link != "" {
		footer = fmt.Sprintf(
			`<p style="margin:24px 0 0;padding-top:16px;border-top:1px solid #e5e7eb;font-size:12px;color:#9ca3af;">`+
				`不想再收到这类邮件？<a href="%s" style="color:#6b7280;">点此退订</a></p>`,
			html.EscapeString(link))
	}

	return fmt.Sprintf(
		`<!DOCTYPE html><html><body style="margin:0;padding:24px;background:#f6f7fb;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;color:#111827;">`+
			`<div style="max-width:560px;margin:0 auto;background:#ffffff;border-radius:12px;padding:24px;">`+
			`<h2 style="margin:0 0 12px;font-size:18px;">%s</h2>`+
			`<p style="margin:0;font-size:14px;line-height:1.8;color:#374151;white-space:pre-wrap;">%s</p>%s`+
			`</div></body></html>`,
		html.EscapeString(title), html.EscapeString(body), footer)
}
