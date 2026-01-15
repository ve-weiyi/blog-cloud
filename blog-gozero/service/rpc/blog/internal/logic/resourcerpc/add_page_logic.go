package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPageLogic {
	return &AddPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建页面
func (l *AddPageLogic) AddPage(in *resourcerpc.AddPageReq) (*resourcerpc.AddPageResp, error) {
	entity := &model.TPage{
		Id:             in.Id,
		PageName:       in.PageName,
		PageLabel:      in.PageLabel,
		PageCover:      in.PageCover,
		IsCarousel:     in.IsCarousel,
		CarouselCovers: jsonconv.AnyToJsonNE(in.CarouselCovers),
	}

	_, err := l.svcCtx.TPageModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.AddPageResp{
		Page: convertPageOut(entity),
	}, nil
}

func convertPageOut(in *model.TPage) (out *resourcerpc.Page) {
	out = &resourcerpc.Page{
		Id:             in.Id,
		PageName:       in.PageName,
		PageLabel:      in.PageLabel,
		PageCover:      in.PageCover,
		IsCarousel:     in.IsCarousel,
		CarouselCovers: nil,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}

	jsonconv.JsonToAny(in.CarouselCovers, &out.CarouselCovers)
	return out
}
