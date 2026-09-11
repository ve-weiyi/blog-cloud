package tag

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/articleservice"
)

type CreateTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewCreateTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTagLogic {
	return &CreateTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTagLogic) CreateTag(req *types.CreateTagReq) (resp *types.TagVO, err error) {
	out, err := l.svcCtx.ArticleService.CreateTag(l.ctx, &articleservice.CreateTagRequest{
		TagName: req.TagName,
	})
	if err != nil {
		return nil, err
	}

	return &types.TagVO{
		Id:           out.Id,
		TagName:      req.TagName,
		ArticleCount: 0,
	}, nil
}
