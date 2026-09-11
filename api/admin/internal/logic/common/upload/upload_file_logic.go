package upload

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/metax"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传文件（服务端上传）
func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, r *http.Request) (resp *types.UploadFileResp, err error) {
	_, h, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}

	f, err := h.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	start := time.Now()
	// 使用 StorageProvider 上传
	fileURL, err := l.svcCtx.StorageProvider.Upload(l.ctx, f, h.Filename)
	if err != nil {
		return nil, err
	}

	// 记录上传日志
	uid, _ := metax.GetApiUserIdFromCtx(l.ctx)
	did, _ := metax.GetApiDeviceIdFromCtx(l.ctx)
	_, _ = l.svcCtx.SyslogService.CreateUploadLog(l.ctx, &syslogservice.CreateUploadLogRequest{
		UserId:   uid,
		DeviceId: did,
		FileBase: req.FileBase,
		FileName: h.Filename,
		FileType: filepath.Ext(h.Filename),
		FileSize: h.Size,
		FileUrl:  fileURL,
	})

	resp = &types.UploadFileResp{
		FileInfo: types.FileInfoVO{
			FileBase:  req.FileBase,
			FileName:  h.Filename,
			FileType:  filepath.Ext(h.Filename),
			FileSize:  h.Size,
			FileUrl:   fileURL,
			UpdatedAt: time.Now().UnixMilli(),
		},
	}

	l.Infof("[Upload] File completed, cost=%dms", time.Since(start).Milliseconds())
	return resp, nil
}
