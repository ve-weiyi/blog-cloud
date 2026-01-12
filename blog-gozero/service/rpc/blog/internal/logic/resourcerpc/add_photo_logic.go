package resourcerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

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
func (l *AddPhotoLogic) AddPhoto(in *resourcerpc.NewPhotoReq) (*resourcerpc.PhotoDetailsResp, error) {
	entity := convertPhotoIn(in)

	_, err := l.svcCtx.TPhotoModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	return convertPhotoOut(entity), nil
}

func convertPhotoIn(in *resourcerpc.NewPhotoReq) (out *model.TPhoto) {
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

func convertPhotoOut(in *model.TPhoto) (out *resourcerpc.PhotoDetailsResp) {
	out = &resourcerpc.PhotoDetailsResp{
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
