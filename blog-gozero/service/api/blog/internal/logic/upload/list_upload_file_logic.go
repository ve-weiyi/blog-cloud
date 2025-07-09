package upload

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取文件列表
func NewListUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUploadFileLogic {
	return &ListUploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUploadFileLogic) ListUploadFile(req *types.ListUploadFileReq) (resp *types.PageResp, err error) {
	up, err := l.svcCtx.Uploader.ListFiles(req.FilePath, int(req.Limit))
	if err != nil {
		return nil, err
	}

	list := make([]*types.FileInfoVO, 0, len(up))
	for _, file := range up {
		vo := &types.FileInfoVO{
			FilePath:  file.FilePath,
			FileName:  file.FileName,
			FileType:  file.FileType,
			FileSize:  file.FileSize,
			FileUrl:   file.FileUrl,
			UpdatedAt: file.UpTime,
		}
		list = append(list, vo)
	}

	resp = &types.PageResp{
		Page:     1,
		PageSize: req.Limit,
		Total:    int64(len(list)),
		List:     list,
	}
	return resp, nil
}
