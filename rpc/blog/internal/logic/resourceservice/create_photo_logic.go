package resourceservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/model"
)

type CreatePhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePhotoLogic {
	return &CreatePhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建照片
func (l *CreatePhotoLogic) CreatePhoto(in *resourcerpc.CreatePhotoRequest) (*resourcerpc.CreatePhotoResponse, error) {
	entity := &model.TPhoto{
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
	}

	_, err := l.svcCtx.TPhotoModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return &resourcerpc.CreatePhotoResponse{Id: entity.Id}, nil
}
