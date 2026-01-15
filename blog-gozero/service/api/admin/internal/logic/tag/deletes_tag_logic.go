package tag

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewDeletesTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesTagLogic {
	return &DeletesTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesTagLogic) DeletesTag(req *types.IdsReq) (resp *types.BatchResp, err error) {
	in := &articlerpc.DeletesTagReq{
		Ids: req.Ids,
	}

	out, err := l.svcCtx.ArticleRpc.DeletesTag(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}
	return resp, nil
}
