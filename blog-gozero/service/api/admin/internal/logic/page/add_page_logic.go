package page

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建页面
func NewAddPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPageLogic {
	return &AddPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPageLogic) AddPage(req *types.NewPageReq) (resp *types.PageBackVO, err error) {
	in := &resourcerpc.AddPageReq{
		Id:             req.Id,
		PageName:       req.PageName,
		PageLabel:      req.PageLabel,
		PageCover:      req.PageCover,
		IsCarousel:     req.IsCarousel,
		CarouselCovers: req.CarouselCovers,
	}

	out, err := l.svcCtx.ResourceRpc.AddPage(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return convertPageTypes(out.Page), nil
}

func convertPageTypes(out *resourcerpc.Page) *types.PageBackVO {
	return &types.PageBackVO{
		Id:             out.Id,
		PageName:       out.PageName,
		PageLabel:      out.PageLabel,
		PageCover:      out.PageCover,
		IsCarousel:     out.IsCarousel,
		CarouselCovers: out.CarouselCovers,
		CreatedAt:      out.CreatedAt,
		UpdatedAt:      out.UpdatedAt,
	}
}
