package event

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertEvent 把通知事件映射为管理端视图。
func convertEvent(in *notificationservice.NotificationEvent) *types.NotificationEventVO {
	if in == nil {
		return nil
	}
	return &types.NotificationEventVO{
		Id:            in.Id,
		Scene:         in.Scene,
		ActorId:       in.ActorId,
		ObjectType:    in.ObjectType,
		ObjectId:      in.ObjectId,
		AudienceType:  in.AudienceType,
		AudienceTotal: in.AudienceTotal,
		Processed:     in.Processed,
		Status:        in.Status,
		Cursor:        in.Cursor,
		CreatedAt:     in.CreatedAt,
		UpdatedAt:     in.UpdatedAt,
	}
}
