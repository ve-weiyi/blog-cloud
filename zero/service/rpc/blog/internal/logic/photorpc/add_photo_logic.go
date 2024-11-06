package photorpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPhotoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPhotoLogic {
	return &AddPhotoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建照片
func (l *AddPhotoLogic) AddPhoto(in *photorpc.PhotoNewReq) (*photorpc.PhotoDetails, error) {
	entity := convertPhotoIn(in)

	_, err := l.svcCtx.TPhotoModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertPhotoOut(entity), nil
}

func convertPhotoIn(in *photorpc.PhotoNewReq) (out *model.TPhoto) {
	out = &model.TPhoto{
		Id:        in.Id,
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
		IsDelete:  in.IsDelete,
	}

	return out
}

func convertPhotoOut(in *model.TPhoto) (out *photorpc.PhotoDetails) {
	out = &photorpc.PhotoDetails{
		Id:        in.Id,
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
		IsDelete:  in.IsDelete,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
	}

	return out
}
