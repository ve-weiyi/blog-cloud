package siteservicelogic

import (
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/x/jsonv"
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
	if in.CarouselCovers != nil {
		if err := jsonv.JsonToAny(*in.CarouselCovers, &out.CarouselCovers); err != nil {
			logx.Errorf("convert page %d carousel covers: %v", in.Id, err)
		}
	}
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
