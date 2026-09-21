package page

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/siteservice"
)

type CreatePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建页面
func NewCreatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePageLogic {
	return &CreatePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePageLogic) CreatePage(req *types.CreatePageReq) (resp *types.PageVO, err error) {
	out, err := l.svcCtx.SiteService.CreatePage(l.ctx, &siteservice.CreatePageRequest{
		PageName:       req.PageName,
		PageLabel:      req.PageLabel,
		PageCover:      req.PageCover,
		IsCarousel:     req.IsCarousel,
		CarouselCovers: req.CarouselCovers,
	})
	if err != nil {
		return nil, err
	}

	return &types.PageVO{
		Id:             out.Id,
		PageName:       req.PageName,
		PageLabel:      req.PageLabel,
		PageCover:      req.PageCover,
		IsCarousel:     req.IsCarousel,
		CarouselCovers: req.CarouselCovers,
	}, nil
}
