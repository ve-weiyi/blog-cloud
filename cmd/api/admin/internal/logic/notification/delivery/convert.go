package delivery

import (
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/notificationservice"
)

// convertDelivery 把投递记录映射为管理端视图。
func convertDelivery(in *notificationservice.Delivery) *types.DeliveryVO {
	if in == nil {
		return nil
	}
	return &types.DeliveryVO{
		Id:             in.Id,
		NotificationId: in.NotificationId,
		Channel:        in.Channel,
		Recipient:      in.Recipient,
		Status:         in.Status,
		TemplateCode:   in.TemplateCode,
		Content:        in.Content,
		ProviderMsgId:  in.ProviderMsgId,
		ErrorMsg:       in.ErrorMsg,
		Attempts:       in.Attempts,
		NextRetryAt:    in.NextRetryAt,
		SentAt:         in.SentAt,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
	}
}
