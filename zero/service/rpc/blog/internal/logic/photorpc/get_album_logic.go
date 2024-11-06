package photorpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlbumLogic {
	return &GetAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取相册
func (l *GetAlbumLogic) GetAlbum(in *photorpc.IdReq) (*photorpc.AlbumDetails, error) {
	entity, err := l.svcCtx.TAlbumModel.FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	cm, err := findPhotoCountGroupAlbum(l.ctx, l.svcCtx, []*model.TAlbum{entity})
	if err != nil {
		return nil, err
	}

	return convertAlbumOut(entity, cm), nil
}
