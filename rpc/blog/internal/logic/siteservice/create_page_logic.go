package siteservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/siterpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
	"github.com/ve-weiyi/vkit/x/jsonconv"
)

type CreatePageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePageLogic {
	return &CreatePageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建页面
func (l *CreatePageLogic) CreatePage(in *siterpc.CreatePageRequest) (*siterpc.CreatePageResponse, error) {
	entity := &model.TPage{
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

	return &siterpc.CreatePageResponse{Id: entity.Id}, nil
}
