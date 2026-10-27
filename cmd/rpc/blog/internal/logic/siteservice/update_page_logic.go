package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/x/jsonv"
)

type UpdatePageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePageLogic {
	return &UpdatePageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新页面
func (l *UpdatePageLogic) UpdatePage(in *siterpc.UpdatePageRequest) (*siterpc.UpdatePageResponse, error) {
	entity, err := l.svcCtx.TPageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.PageName = in.PageName
	entity.PageLabel = in.PageLabel
	entity.PageCover = in.PageCover
	entity.IsCarousel = in.IsCarousel
	covers := jsonv.AnyToJsonNE(in.CarouselCovers)
	entity.CarouselCovers = &covers

	_, err = l.svcCtx.TPageModel.Update(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &siterpc.UpdatePageResponse{
		Page: convertPageOut(entity),
	}, nil
}
