package tag

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/contentservice"
)

type BatchDeleteTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除标签
func NewBatchDeleteTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteTagLogic {
	return &BatchDeleteTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteTagLogic) BatchDeleteTag(req *types.DeleteTagReq) (resp *types.BatchResp, err error) {
	out, err := l.svcCtx.ContentService.BatchDeleteTags(l.ctx, &contentservice.BatchDeleteTagsRequest{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BatchResp{
		SuccessCount: out.SuccessCount,
	}, nil
}
