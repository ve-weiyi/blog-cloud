package upload

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
)

type QueryFileListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询上传文件列表
func NewQueryFileListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryFileListLogic {
	return &QueryFileListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryFileListLogic) QueryFileList(req *types.QueryFileListReq) (resp *types.QueryFileListResp, err error) {
	files, err := l.svcCtx.StorageProvider.ListFiles(l.ctx, req.FileBase, 100)
	if err != nil {
		return nil, err
	}

	var filtered []*types.FileInfoVO
	for _, f := range files {
		if req.Keyword != "" && !strings.Contains(f.FileName, req.Keyword) {
			continue
		}
		filtered = append(filtered, &types.FileInfoVO{
			FileBase:  f.FilePath,
			FileName:  f.FileName,
			FileType:  f.FileType,
			FileSize:  f.FileSize,
			FileUrl:   f.FileURL,
			UpdatedAt: f.UpTime.UnixMilli(),
		})
	}

	total := int64(len(filtered))
	start := (req.Page - 1) * req.PageSize
	end := start + req.PageSize
	if start < 0 {
		start = 0
	}
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return &types.QueryFileListResp{
		List:     filtered[start:end],
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
