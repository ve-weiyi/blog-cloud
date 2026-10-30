package user_notification

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertItem 把条目级视图映射为管理端视图。
func convertItem(in *notificationservice.NotificationItem) *types.NotificationItemVO {
	if in == nil {
		return nil
	}
	return &types.NotificationItemVO{
		Id:         in.Id,
		EventId:    in.EventId,
		Scene:      in.Scene,
		ActorId:    in.ActorId,
		ActorName:  in.ActorName,
		ObjectType: in.ObjectType,
		ObjectId:   in.ObjectId,
		Title:      in.Title,
		Content:    in.Content,
		Link:       in.Link,
		ReadAt:     in.ReadAt,
		CreatedAt:  in.CreatedAt,
	}
}
