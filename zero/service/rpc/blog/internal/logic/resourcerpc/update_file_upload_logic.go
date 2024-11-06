package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileUploadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFileUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileUploadLogic {
	return &UpdateFileUploadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新文件上传
func (l *UpdateFileUploadLogic) UpdateFileUpload(in *resourcerpc.FileUploadNewReq) (*resourcerpc.FileUploadDetails, error) {
	entity := &model.TFileUpload{
		Id:       in.Id,
		UserId:   in.UserId,
		FilePath: in.FilePath,
		FileType: in.FileType,
		FileName: in.FileName,
		FileSize: in.FileSize,
		FileMd5:  in.FileMd5,
		FileUrl:  in.FileUrl,
	}

	_, err := l.svcCtx.TFileUploadModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFileUploadOut(entity), nil
}
