package contentservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/contentrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type BatchDeleteTagsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDeleteTagsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteTagsLogic {
	return &BatchDeleteTagsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量删除标签
func (l *BatchDeleteTagsLogic) BatchDeleteTags(in *contentrpc.BatchDeleteTagsRequest) (*contentrpc.BatchDeleteTagsResponse, error) {
	rows, err := l.svcCtx.TTagModel.DeleteBatch(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &contentrpc.BatchDeleteTagsResponse{SuccessCount: rows}, nil
}
