package siteservicelogic

import (
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"
)

func convertPageOut(in *model.TPage) *siterpc.Page {
	out := &siterpc.Page{
		Id:         in.Id,
		PageName:   in.PageName,
		PageLabel:  in.PageLabel,
		PageCover:  in.PageCover,
		IsCarousel: in.IsCarousel,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
	}
	jsonconv.JsonToAny(in.CarouselCovers, &out.CarouselCovers)
	return out
}

func convertFriendOut(in *model.TFriend) *siterpc.Friend {
	return &siterpc.Friend{
		Id:          in.Id,
		LinkName:    in.LinkName,
		LinkAvatar:  in.LinkAvatar,
		LinkAddress: in.LinkAddress,
		LinkIntro:   in.LinkIntro,
		CreatedAt:   in.CreatedAt.UnixMilli(),
		UpdatedAt:   in.UpdatedAt.UnixMilli(),
	}
}
