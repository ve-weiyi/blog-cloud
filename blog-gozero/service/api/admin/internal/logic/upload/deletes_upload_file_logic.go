package upload

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletesUploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除文件列表
func NewDeletesUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletesUploadFileLogic {
	return &DeletesUploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletesUploadFileLogic) DeletesUploadFile(req *types.DeletesUploadFileReq) (resp *types.BatchResp, err error) {
	for _, filepath := range req.FilePaths {
		err = l.svcCtx.Uploader.DeleteFile(filepath)
		if err != nil {
			return nil, err
		}
	}

	return &types.BatchResp{
		SuccessCount: int64(len(req.FilePaths)),
	}, nil
}
