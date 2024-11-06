package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddFileFolderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddFileFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFileFolderLogic {
	return &AddFileFolderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建文件夹
func (l *AddFileFolderLogic) AddFileFolder(in *resourcerpc.FileFolderNewReq) (*resourcerpc.FileFolderDetails, error) {

	entity := &model.TFileFolder{
		Id:         0,
		UserId:     in.UserId,
		FilePath:   in.FilePath,
		FolderName: in.FolderName,
		FolderDesc: in.FolderDesc,
	}

	_, err := l.svcCtx.TFileFolderModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFileFolderOut(entity), nil
}

func convertFileFolderOut(in *model.TFileFolder) (out *resourcerpc.FileFolderDetails) {
	out = &resourcerpc.FileFolderDetails{
		Id:         in.Id,
		UserId:     in.UserId,
		FilePath:   in.FilePath,
		FolderName: in.FolderName,
		FolderDesc: in.FolderDesc,
		CreatedAt:  in.CreatedAt.Unix(),
		UpdatedAt:  in.UpdatedAt.Unix(),
	}

	return out
}
