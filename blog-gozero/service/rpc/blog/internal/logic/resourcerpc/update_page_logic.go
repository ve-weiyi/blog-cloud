package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *UpdatePageLogic) UpdatePage(in *resourcerpc.UpdatePageReq) (*resourcerpc.UpdatePageResp, error) {
	entity, err := l.svcCtx.TPageModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	entity.PageName = in.PageName
	entity.PageLabel = in.PageLabel
	entity.PageCover = in.PageCover
	entity.IsCarousel = in.IsCarousel
	entity.CarouselCovers = jsonconv.AnyToJsonNE(in.CarouselCovers)

	_, err = l.svcCtx.TPageModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.UpdatePageResp{
		Page: convertPageOut(entity),
	}, nil
}
