package notify

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertGroup 把收件箱分组映射为读者端视图。
//
// 字段与管理端的同名结构保持一致：两个网关的生成类型互不可见，只能各存一份。
func convertGroup(in *notificationservice.NotificationGroup) *types.NotificationGroupVO {
	if in == nil {
		return nil
	}
	return &types.NotificationGroupVO{
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
