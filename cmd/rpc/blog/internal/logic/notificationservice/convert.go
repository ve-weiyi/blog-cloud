package notificationservicelogic

import (
	"time"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/notification"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/model"
)

func unix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.Unix()
}

func unixOf(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func convertAnnouncement(in *model.TAnnouncement) *notificationrpc.Announcement {
	if in == nil {
		return nil
	}
	return &notificationrpc.Announcement{
		Id:          in.Id,
		Title:       in.Title,
		Content:     deref(in.Content),
		Category:    in.Category,
		Level:       in.Level,
		TargetType:  in.TargetType,
		TargetIds:   deref(in.TargetIds),
		Status:      in.Status,
		PublishedAt: unixOf(in.PublishedAt),
		PublishedBy: in.PublishedBy,
		CreatedAt:   unix(in.CreatedAt),
		UpdatedAt:   unix(in.UpdatedAt),
	}
}

func convertDelivery(in *model.TNotificationDelivery) *notificationrpc.Delivery {
	if in == nil {
		return nil
	}
	return &notificationrpc.Delivery{
		Id:             in.Id,
		NotificationId: in.NotificationId,
		Channel:        in.Channel,
		Recipient:      in.Recipient,
		Status:         in.Status,
		TemplateCode:   in.TemplateCode,
		Content:        deref(in.Content),
		ProviderMsgId:  in.ProviderMsgId,
		ErrorMsg:       in.ErrorMsg,
		Attempts:       in.Attempts,
		NextRetryAt:    unixOf(in.NextRetryAt),
		SentAt:         unixOf(in.SentAt),
		CreatedAt:      unix(in.CreatedAt),
		UpdatedAt:      unix(in.UpdatedAt),
	}
}

func convertEvent(in *model.TNotificationEvent) *notificationrpc.NotificationEvent {
	if in == nil {
		return nil
	}
	return &notificationrpc.NotificationEvent{
		Id:            in.Id,
		Scene:         in.Scene,
		ActorId:       in.ActorId,
		ObjectType:    in.ObjectType,
		ObjectId:      in.ObjectId,
		AudienceType:  in.AudienceType,
		Status:        in.Status,
		AudienceTotal: in.AudienceTotal,
		Processed:     in.Processed,
		Cursor:        in.Cursor,
		CreatedAt:     unix(in.CreatedAt),
		UpdatedAt:     unix(in.UpdatedAt),
	}
}

func convertInboxGroup(in notification.InboxGroup) *notificationrpc.NotificationGroup {
	return &notificationrpc.NotificationGroup{
		GroupKey:        in.GroupKey,
		Scene:           in.Scene,
		Title:           in.Title,
		Content:         in.Content,
		Link:            in.Link,
		Level:           in.Level,
		ActorCount:      in.ActorCount,
		LatestActorName: in.LatestActorName,
		LatestAt:        in.LatestAt,
		Unread:          in.Unread,
	}
}

func convertItem(in notification.Item) *notificationrpc.NotificationItem {
	return &notificationrpc.NotificationItem{
		Id:          in.ID,
		EventId:     in.EventID,
		RecipientId: in.RecipientID,
		Scene:       in.Scene,
		ActorId:     in.ActorID,
		ActorName:   in.ActorName,
		ObjectType:  in.ObjectType,
		ObjectId:    in.ObjectID,
		Title:       in.Title,
		Content:     in.Content,
		Link:        in.Link,
		ReadAt:      unixOf(in.ReadAt),
		CreatedAt:   unix(in.CreatedAt),
	}
}

func convertPreference(in notification.Preference) *notificationrpc.NotificationPreference {
	return &notificationrpc.NotificationPreference{
		Group:     string(in.Group),
		Channel:   string(in.Channel),
		Enabled:   in.Enabled,
		Mandatory: in.Mandatory,
	}
}
