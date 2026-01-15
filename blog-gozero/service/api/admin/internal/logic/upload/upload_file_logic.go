package upload

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/kit/oss"
	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/cryptox"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 上传文件
func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, r *http.Request) (resp *types.FileInfoVO, err error) {
	f, h, err := r.FormFile("file")
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

	return &types.FileInfoVO{
		FilePath:  req.FilePath,
		FileName:  out.FileLog.FileName,
		FileType:  out.FileLog.FileType,
		FileSize:  out.FileLog.FileSize,
		FileUrl:   out.FileLog.FileUrl,
		UpdatedAt: out.FileLog.UpdatedAt,
	}, nil
}
