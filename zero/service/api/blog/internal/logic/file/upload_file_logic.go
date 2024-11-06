package file

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/crypto"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/resourcerpc"

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

func (l *UploadFileLogic) UploadFile(req *types.UploadFileReq, r *http.Request) (resp *types.FileBackDTO, err error) {
	f, h, _ := r.FormFile("file")
	defer f.Close()

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

	return ConvertFileUploadTypes(out), nil
}

func ConvertFileUploadTypes(in *resourcerpc.FileUploadDetails) (out *types.FileBackDTO) {
	out = &types.FileBackDTO{
		Id:        in.Id,
		UserId:    in.UserId,
		FilePath:  in.FilePath,
		FileName:  in.FileName,
		FileType:  in.FileType,
		FileSize:  in.FileSize,
		FileMd5:   in.FileMd5,
		FileUrl:   in.FileUrl,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
