package mediaservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/mediarpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/model"
)

type CreateAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAlbumLogic {
	return &CreateAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建相册
func (l *CreateAlbumLogic) CreateAlbum(in *mediarpc.CreateAlbumRequest) (*mediarpc.CreateAlbumResponse, error) {
	if !enums.IsValidAlbumStatus(in.Status) {
		return nil, bizerr.NewBizError(bizcode.CodeParamValueInvalid, "相册状态不合法")
	}

	entity := &model.TAlbum{
		AlbumName:  in.AlbumName,
		AlbumDesc:  in.AlbumDesc,
		AlbumCover: in.AlbumCover,
		Status:     in.Status,
	}

	_, err := l.svcCtx.TAlbumModel.Insert(l.ctx, entity)
	if err != nil {
		return nil, err
	}

	// 新相册还没有照片，计数为空
	return &mediarpc.CreateAlbumResponse{Id: entity.Id, Album: convertAlbumOut(entity, map[int64]int{})}, nil
}
