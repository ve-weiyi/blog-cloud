package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileFolderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileFolderLogic {
	return &DeleteFileFolderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文件夹
func (l *DeleteFileFolderLogic) DeleteFileFolder(in *resourcerpc.IdsReq) (*resourcerpc.BatchResp, error) {
	rows, err := l.svcCtx.TFileFolderModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
