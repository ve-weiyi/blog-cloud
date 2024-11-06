package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileUploadLogic {
	return &DeleteFileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文件上传
func (l *DeleteFileUploadLogic) DeleteFileUpload(in *resourcerpc.IdsReq) (*resourcerpc.BatchResp, error) {
	rows, err := l.svcCtx.TFileUploadModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
