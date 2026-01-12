package page

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新页面
func NewUpdatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePageLogic {
	return &UpdatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePageLogic) UpdatePage(req *types.NewPageReq) (resp *types.PageBackVO, err error) {
	in := &resourcerpc.NewPageReq{
		Id:             req.Id,
		PageName:       req.PageName,
		PageLabel:      req.PageLabel,
		PageCover:      req.PageCover,
		IsCarousel:     req.IsCarousel,
		CarouselCovers: req.CarouselCovers,
	}

	out, err := l.svcCtx.ResourceRpc.UpdatePage(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.PageBackVO{
		Id:             out.Id,
		PageName:       out.PageName,
		PageLabel:      out.PageLabel,
		PageCover:      out.PageCover,
		IsCarousel:     out.IsCarousel,
		CarouselCovers: out.CarouselCovers,
		CreatedAt:      out.CreatedAt,
		UpdatedAt:      out.UpdatedAt,
	}
	return resp, nil
}
