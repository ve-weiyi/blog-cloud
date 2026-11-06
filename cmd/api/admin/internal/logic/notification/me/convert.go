package me

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertGroup 把收件箱分组映射为管理端视图。
//
// 与读者端同名结构保持字段一致：两个网关的生成类型互不可见，只能各存一份。
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
