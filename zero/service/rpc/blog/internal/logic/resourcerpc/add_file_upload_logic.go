package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFileUploadLogic {
	return &AddFileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文件上传
func (l *AddFileUploadLogic) AddFileUpload(in *resourcerpc.FileUploadNewReq) (*resourcerpc.FileUploadDetails, error) {
	entity := &model.TFileUpload{
		Id:       0,
		UserId:   in.UserId,
		FilePath: in.FilePath,
		FileName: in.FileName,
		FileType: in.FileType,
		FileSize: in.FileSize,
		FileMd5:  in.FileMd5,
		FileUrl:  in.FileUrl,
	}

	_, err := l.svcCtx.TFileUploadModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFileUploadOut(entity), nil
}

func convertFileUploadOut(in *model.TFileUpload) (out *resourcerpc.FileUploadDetails) {
	out = &resourcerpc.FileUploadDetails{
		Id:        in.Id,
		UserId:    in.UserId,
		FilePath:  in.FilePath,
		FileName:  in.FileName,
		FileType:  in.FileType,
		FileSize:  in.FileSize,
		FileMd5:   in.FileMd5,
		FileUrl:   in.FileUrl,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
	}

	return out
}
