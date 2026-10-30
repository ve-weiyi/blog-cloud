package notification

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
	"github.com/ve-weiyi/blog-cloud/model"
	mqx "github.com/ve-weiyi/vkit/adapter/mqx"
)

// Enqueuer 把唤醒消息写入发件箱。与业务写入同事务，保证「库里有事件，消息才会发」。
type Enqueuer interface {
	Enqueue(ctx context.Context, tx *gorm.DB, msg mqx.Message) error
}

// Notifier 在条目写入后发出实时信号。为 nil 时跳过 —— 关掉实时通道，功能仍然正确。
type Notifier interface {
	NotifyChanged(ctx context.Context, recipients []string)
}

// AnnouncementNotifier 广播公告状态变化，供管理端列表刷新。
// 与读者端的「收件箱变了」是两种语义，频道与载荷分开。
type AnnouncementNotifier interface {
	NotifyAnnouncement(ctx context.Context, published bool, id int64, title, category, level string, at time.Time)
}

// 事件状态机：待扇出 → 扇出中 → 已完成 / 已失败 / 已终止。
//
// 只有已终止且已无引用条目的事件才允许清理；仍在扇出中的事件被清掉，
// 扇出会断在半路且没有任何痕迹。
const (
	StatusPending    = "pending"
	StatusFanning    = "fanning"
	StatusDone       = "done"
	StatusFailed     = "failed"
	StatusTerminated = "terminated"
)

// 投递状态。已取消用于「发送前复查发现已退订或条目已删」，置为终态才不会
// 被重试扫描反复捞起。
const (
	DeliveryPending   = "pending"
	DeliverySent      = "sent"
	DeliveryFailed    = "failed"
	DeliveryCancelled = "cancelled"
)

// 公告状态。
const (
	AnnouncementDraft     = "draft"
	AnnouncementPublished = "published"
	AnnouncementRevoked   = "revoked"
)

// IsTerminal 判断事件是否已经走到终态。
func IsTerminal(status string) bool {
	return status == StatusDone || status == StatusFailed || status == StatusTerminated
}

// Deps 是通知服务的依赖。依赖项较多，用结构体而非位置参数：位置错一个只会静默错位。
type Deps struct {
	DB            *gorm.DB
	Events        model.TNotificationEventModel
	Notifications model.TNotificationModel
	Deliveries    model.TNotificationDeliveryModel
	Preferences   model.TNotificationPreferenceModel
	Announcements model.TAnnouncementModel
	Users         model.TUserModel
	Comments      model.TCommentModel
	Outbox        Enqueuer
	Notifier      Notifier // 可为 nil：不注入时只写库，不发实时信号

	AnnouncementNotifier AnnouncementNotifier // 可为 nil：不注入时管理端列表只能手动刷新

	// UnsubscribeSecret / SiteURL 用于在邮件正文里签发退订链接；
	// 任一为空时不渲染链接（拼一个坏链接比没有链接更糟）。
	UnsubscribeSecret string
	SiteURL           string
}

// Service 是通知域的实现。依赖只在构造时注入，方法不再接收依赖。
type Service struct {
	db            *gorm.DB
	events        model.TNotificationEventModel
	notifs        model.TNotificationModel
	deliveries    model.TNotificationDeliveryModel
	prefs         model.TNotificationPreferenceModel
	announcements model.TAnnouncementModel
	users         model.TUserModel
	comments      model.TCommentModel
	outbox        Enqueuer
	notifier      Notifier

	announcementNotifier AnnouncementNotifier

	unsubscribeSecret string
	siteURL           string
}

// ResolveActorName 取触发者展示名：注册用户取昵称。
//
// 取不到时返回空串，由渲染层退化为「有人」；访客没有昵称，调用方传「访客」更贴切。
func (s *Service) ResolveActorName(ctx context.Context, userID string) string {
	if userID == "" {
		return ""
	}
	row, err := s.users.FindOneByUserId(ctx, userID)
	if err != nil || row.Nickname == "" {
		return ""
	}
	return row.Nickname
}

// NewService 构造通知服务。
func NewService(d Deps) *Service {
	// 不阻断启动：缺密钥时退订链接不签发、退订接口一律拒绝，功能退化但不产生坏链接。
	// 但必须让运维看见 —— 否则线上静默缺失退订能力，直到有人点邮件才发现
	if len(d.UnsubscribeSecret) < notify.MinSecretLen || d.SiteURL == "" {
		logx.Errorf("通知配置未生效：unsubscribe_secret 少于 %d 字节或 site_url 为空，"+
			"邮件里的退订链接不会签发", notify.MinSecretLen)
	}

	return &Service{
		db:                   d.DB,
		events:               d.Events,
		notifs:               d.Notifications,
		deliveries:           d.Deliveries,
		prefs:                d.Preferences,
		announcements:        d.Announcements,
		users:                d.Users,
		comments:             d.Comments,
		outbox:               d.Outbox,
		notifier:             d.Notifier,
		announcementNotifier: d.AnnouncementNotifier,
		unsubscribeSecret:    d.UnsubscribeSecret,
		siteURL:              d.SiteURL,
	}
}
