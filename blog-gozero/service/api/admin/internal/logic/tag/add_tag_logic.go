package tag

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建标签
func NewAddTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagLogic {
	return &AddTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddTagLogic) AddTag(req *types.NewTagReq) (resp *types.TagBackVO, err error) {
	in := &articlerpc.NewTagReq{
		Id:      req.Id,
		TagName: req.TagName,
	}

	out, err := l.svcCtx.ArticleRpc.AddTag(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.TagBackVO{
		Id:           out.Id,
		TagName:      out.TagName,
		ArticleCount: 0,
		CreatedAt:    out.CreatedAt,
		UpdatedAt:    out.UpdatedAt,
	}
	return resp, nil
}
