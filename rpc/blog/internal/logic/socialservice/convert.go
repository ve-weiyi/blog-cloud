package socialservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"

	"github.com/ve-weiyi/blog-cloud/service/app/rpc/internal/pb/socialrpc"
)

func convertFriendOut(in *model.TFriend) *socialrpc.Friend {
	return &socialrpc.Friend{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		CreatedAt:   in.CreatedAt.UnixMilli(),
		UpdatedAt:   in.UpdatedAt.UnixMilli(),
	}
}

func convertTalkOut(in *model.TTalk) *socialrpc.Talk {
	var images []string
	jsonconv.JsonToAny(in.Images, &images)

	return &socialrpc.Talk{
		Id:        in.Id,
		UserId:    in.UserId,
		Content:   in.Content,
		Images:    images,
		IsTop:     in.IsTop,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
		LikeCount: in.LikeCount,
	}
}
