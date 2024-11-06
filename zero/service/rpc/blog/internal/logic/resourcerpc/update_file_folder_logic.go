package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileFolderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFileFolderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileFolderLogic {
	return &UpdateFileFolderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新文件夹
func (l *UpdateFileFolderLogic) UpdateFileFolder(in *resourcerpc.FileFolderNewReq) (*resourcerpc.FileFolderDetails, error) {
	entity := &model.TFileFolder{
		Id:         in.Id,
		UserId:     in.UserId,
		FilePath:   in.FilePath,
		FolderName: in.FolderName,
		FolderDesc: in.FolderDesc,
	}

	_, err := l.svcCtx.TFileFolderModel.Save(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertFileFolderOut(entity), nil
}
