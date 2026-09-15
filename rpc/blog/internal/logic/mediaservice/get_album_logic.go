package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
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

// 查询相册详情
func (l *GetAlbumLogic) GetAlbum(in *mediarpc.GetAlbumRequest) (*mediarpc.GetAlbumResponse, error) {
	entity, err := l.svcCtx.TAlbumModel.FindById(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &mediarpc.GetAlbumResponse{
		Album: convertAlbumOut(entity, nil),
	}, nil
}
