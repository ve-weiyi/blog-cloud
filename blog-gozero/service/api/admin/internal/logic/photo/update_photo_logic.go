package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePhotoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新照片
func NewUpdatePhotoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePhotoLogic {
	return &UpdatePhotoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePhotoLogic) UpdatePhoto(req *types.PhotoNewReq) (resp *types.PhotoBackVO, err error) {
	in := &resourcerpc.PhotoNewReq{
		Id:        req.Id,
		AlbumId:   req.AlbumId,
		PhotoName: req.PhotoName,
		PhotoDesc: req.PhotoDesc,
		PhotoSrc:  req.PhotoSrc,
		IsDelete:  req.IsDelete,
	}

	out, err := l.svcCtx.ResourceRpc.UpdatePhoto(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.PhotoBackVO{
		Id:        out.Id,
		AlbumId:   out.AlbumId,
		PhotoName: out.PhotoName,
		PhotoDesc: out.PhotoDesc,
		PhotoSrc:  out.PhotoSrc,
		IsDelete:  out.IsDelete,
		CreatedAt: out.CreatedAt,
		UpdatedAt: out.UpdatedAt,
	}
	return resp, nil
}
