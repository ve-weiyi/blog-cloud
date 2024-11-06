package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddPhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建照片
func NewAddPhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPhotoLogic {
	return &AddPhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPhotoLogic) AddPhoto(req *types.PhotoNewReq) (resp *types.PhotoBackDTO, err error) {
	in := ConvertPhotoPb(req)
	out, err := l.svcCtx.PhotoRpc.AddPhoto(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertPhotoTypes(out)
	return resp, nil
}

func ConvertPhotoPb(in *types.PhotoNewReq) (out *photorpc.PhotoNewReq) {
	out = &photorpc.PhotoNewReq{
		Id:        in.Id,
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
		IsDelete:  in.IsDelete,
	}

	return
}

func ConvertPhotoTypes(in *photorpc.PhotoDetails) (out *types.PhotoBackDTO) {
	out = &types.PhotoBackDTO{
		Id:        in.Id,
		AlbumId:   in.AlbumId,
		PhotoName: in.PhotoName,
		PhotoDesc: in.PhotoDesc,
		PhotoSrc:  in.PhotoSrc,
		IsDelete:  in.IsDelete,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}

	return
}
