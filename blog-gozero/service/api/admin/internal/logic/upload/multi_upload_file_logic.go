package upload

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oss"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"
)

type MultiUploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传文件列表
func NewMultiUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUploadFileLogic {
	return &MultiUploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUploadFileLogic) MultiUploadFile(req *types.MultiUploadFileReq, r *http.Request) (resp []*types.FileInfoVO, err error) {
	start := time.Now()
	files := r.MultipartForm.File["files"]

	for i, h := range files {
		fileStart := time.Now()

		f, err := h.Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()

		up, err := l.svcCtx.Uploader.UploadFile(f, req.FilePath, oss.NewFileNameWithDateTime(h.Filename))
		if err != nil {
			return nil, err
		}

		in := &syslogrpc.AddFileLogReq{
			FilePath: req.FilePath,
			FileName: h.Filename,
			FileType: filepath.Ext(h.Filename),
			FileSize: h.Size,
			FileMd5:  cryptox.Md5v(h.Filename, ""),
			FileUrl:  up,
		}

		rpcCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		out, err := l.svcCtx.SyslogRpc.AddFileLog(rpcCtx, in)
		cancel()
		if err != nil {
			return nil, err
		}

		resp = append(resp, &types.FileInfoVO{
			FilePath:  req.FilePath,
			FileName:  out.FileLog.FileName,
			FileType:  out.FileLog.FileType,
			FileSize:  out.FileLog.FileSize,
			FileUrl:   out.FileLog.FileUrl,
			UpdatedAt: out.FileLog.UpdatedAt,
		})

		l.Infof("[Upload] File %d/%d completed, cost=%dms", i+1, len(files), time.Since(fileStart).Milliseconds())
	}

	l.Infof("[Upload] All files uploaded, total=%d, cost=%dms", len(files), time.Since(start).Milliseconds())
	return resp, nil
}
