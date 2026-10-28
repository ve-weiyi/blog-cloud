package consumer

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

// channel 是邮件与短信在投递流程中的差异点。
//
// resolve 可以在解析模板时改写自己的事件，send 随后从事件读取渠道私有的字段
// （如邮件的标题）。
type channel interface {
	// name 返回渠道名，用于投递记录。
	name() string
	// recipient 返回收件人。
	recipient() string
	// resolve 解析模板，返回模板编码与正文。
	resolve(ctx context.Context) (templateCode, content string, err error)
	// send 执行投递。
	send(ctx context.Context, templateCode, content string) error
	// retryable 判断发送失败是否值得重试。
	retryable(err error) bool
}

// delivery 是邮件与短信共用的投递编排：按事件 ID 幂等查重，
// 创建或复位投递记录，发送后回写结果。
type delivery struct {
	records model.TNotifyRecordModel
	ch      channel
	// bizID 是事件的幂等键，由消费入口从投递上下文取出。
	bizID string
}

// run 执行一次投递。返回非 nil 表示本次失败且值得重试。
func (d *delivery) run(ctx context.Context) error {
	logger := logx.WithContext(ctx)

	if d.bizID == "" {
		return mqx.Permanent(errors.New("mq: event id is missing"))
	}

	existing, err := findDelivery(ctx, d.records, d.bizID, d.ch.name())
	if err != nil {
		return err
	}
	if existing != nil && existing.Status == enums.NotifyRecordStatusSent {
		logger.Infof("%s 已发送，跳过重复投递: bizId=%s", d.ch.name(), d.bizID)
		return nil
	}

	templateCode, content, err := d.ch.resolve(ctx)
	if err != nil {
		return err
	}

	record, err := d.start(ctx, existing, d.bizID, templateCode, content)
	if err != nil || record == nil {
		return err
	}

	return d.finish(ctx, record, d.ch.send(ctx, templateCode, content))
}

// start 创建或复位投递记录；返回 nil 表示并发消费已把它置为已发送。
func (d *delivery) start(ctx context.Context, existing *model.TNotifyRecord, bizID, templateCode, content string) (*model.TNotifyRecord, error) {
	logger := logx.WithContext(ctx)

	if existing == nil {
		record, err := createDelivery(ctx, d.records, &model.TNotifyRecord{
			Channel:      d.ch.name(),
			Recipient:    d.ch.recipient(),
			TemplateCode: templateCode,
			Content:      stringToPtr(content),
			Status:       enums.NotifyRecordStatusPending,
			BizId:        bizID,
			CreatedAt:    time.Now(),
		})
		if err != nil {
			logger.Errorf("插入 %s 投递记录失败: %v", d.ch.name(), err)
			return nil, err
		}
		if record.Status == enums.NotifyRecordStatusSent {
			logger.Infof("%s 已发送，跳过重复投递: bizId=%s", d.ch.name(), bizID)
			return nil, nil
		}
		return record, nil
	}

	fields := map[string]interface{}{
		"recipient":     d.ch.recipient(),
		"template_code": templateCode,
		"content":       content,
		"status":        enums.NotifyRecordStatusPending,
		"error_msg":     "",
	}
	if _, err := d.records.UpdateFields(ctx, fields, "id = ?", existing.Id); err != nil {
		logger.Errorf("更新 %s 投递记录失败: %v", d.ch.name(), err)
		return nil, err
	}
	return existing, nil
}

// finish 回写发送结果。不可重试的失败被记为已处理，消息直接确认。
func (d *delivery) finish(ctx context.Context, record *model.TNotifyRecord, sendErr error) error {
	logger := logx.WithContext(ctx)
	fields := make(map[string]interface{})

	if sendErr != nil {
		logger.Errorf("发送%s失败: %v", d.ch.name(), sendErr)
		fields["status"] = enums.NotifyRecordStatusFailed
		fields["error_msg"] = sendErr.Error()

		if !d.ch.retryable(sendErr) {
			logger.Infof("不可重试的错误，直接确认消息: %v", sendErr)
			_, _ = d.records.UpdateFields(ctx, fields, "id = ?", record.Id)
			return nil
		}
	} else {
		logger.Infof("发送%s成功: %s", d.ch.name(), record.Recipient)
		fields["status"] = enums.NotifyRecordStatusSent
		fields["sent_at"] = time.Now()
	}

	if _, err := d.records.UpdateFields(ctx, fields, "id = ?", record.Id); err != nil {
		logger.Errorf("更新%s状态失败: %v", d.ch.name(), err)
	}

	return sendErr
}

// findDelivery 按幂等键查找投递记录；不存在返回 nil, nil。
func findDelivery(ctx context.Context, m model.TNotifyRecordModel, bizID, channel string) (*model.TNotifyRecord, error) {
	if bizID == "" {
		return nil, nil
	}

	record, err := m.FindOne(ctx, "", "biz_id = ? AND channel = ?", bizID, channel)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

// createDelivery 插入投递记录；并发下唯一键冲突时返回已存在的记录。
func createDelivery(ctx context.Context, m model.TNotifyRecordModel, record *model.TNotifyRecord) (*model.TNotifyRecord, error) {
	if _, err := m.Insert(ctx, record); err != nil {
		if !errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, err
		}

		existing, findErr := findDelivery(ctx, m, record.BizId, record.Channel)
		if findErr != nil {
			return nil, findErr
		}
		if existing == nil {
			return nil, err
		}
		return existing, nil
	}
	return record, nil
}
