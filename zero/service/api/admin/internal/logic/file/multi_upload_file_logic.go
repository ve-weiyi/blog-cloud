package file

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *MultiUploadFileLogic) MultiUploadFile(req *types.MultiUploadFileReq, r *http.Request) (resp []*types.FileBackDTO, err error) {
	// 获取文件切片
	files := r.MultipartForm.File["files"]
	for _, h := range files {
		up, err := l.svcCtx.Uploader.UploadFile(req.FilePath, h)
		if err != nil {
			return nil, err
		}

		in := &resourcerpc.FileUploadNewReq{
			UserId:   cast.ToString(l.ctx.Value("uid")),
			FilePath: req.FilePath,
			FileName: h.Filename,
			FileType: filepath.Ext(h.Filename),
			FileSize: h.Size,
			FileMd5:  crypto.Md5v(h.Filename, ""),
			FileUrl:  up,
		}

		out, err := l.svcCtx.ResourceRpc.AddFileUpload(l.ctx, in)
		if err != nil {
			return nil, err
		}

		resp = append(resp, ConvertFileUploadTypes(out))
	}

	return
}
