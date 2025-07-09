package articlerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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

// 查询标签
func (l *GetTagLogic) GetTag(in *articlerpc.IdReq) (*articlerpc.TagDetailsResp, error) {
	entity, err := l.svcCtx.TTagModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &articlerpc.TagDetailsResp{
		Id:           entity.Id,
		TagName:      entity.TagName,
		CreatedAt:    entity.CreatedAt.Unix(),
		UpdatedAt:    entity.UpdatedAt.Unix(),
		ArticleCount: 0,
	}, nil
}
