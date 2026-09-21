package chatservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/chatrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

func convertChatOut(in *model.TChat) *chatrpc.Chat {
	return &chatrpc.Chat{
		Id:        in.Id,
		UserId:    in.UserId,
		DeviceId:  in.DeviceId,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Type:      in.Type,
		Content:   in.Content,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
	}
}
