package announcement

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertAnnouncement 把公告实体映射为管理端视图。
func convertAnnouncement(in *notificationservice.Announcement) *types.AnnouncementVO {
	if in == nil {
		return nil
	}
	return &types.AnnouncementVO{
		Id:          in.Id,
		Title:       in.Title,
		Content:     in.Content,
		Category:    in.Category,
		Level:       in.Level,
		TargetType:  in.TargetType,
		TargetIds:   in.TargetIds,
		Status:      in.Status,
		PublishedAt: in.PublishedAt,
		PublishedBy: in.PublishedBy,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}
