package photorpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/photorpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAlbumLogic {
	return &DeleteAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除相册
func (l *DeleteAlbumLogic) DeleteAlbum(in *photorpc.IdsReq) (*photorpc.BatchResp, error) {
	rows, err := l.svcCtx.TAlbumModel.Deletes(l.ctx, "id in (?)", in.Ids)
	if err != nil {
		return nil, err
	}

	return &photorpc.BatchResp{
		SuccessCount: rows,
	}, nil
}
