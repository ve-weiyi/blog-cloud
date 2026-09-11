package articleservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagLogic {
	return &GetTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询标签详情
func (l *GetTagLogic) GetTag(in *articlerpc.GetTagRequest) (*articlerpc.GetTagResponse, error) {
	entity, err := l.svcCtx.TTagModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &articlerpc.GetTagResponse{
		Tag: &articlerpc.Tag{
			Id:        entity.Id,
			TagName:   entity.TagName,
			CreatedAt: entity.CreatedAt.UnixMilli(),
			UpdatedAt: entity.UpdatedAt.UnixMilli(),
		},
	}, nil
}
