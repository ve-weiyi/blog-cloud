package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"
)

type FindPhotoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取照片列表
func NewFindPhotoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindPhotoListLogic {
	return &FindPhotoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindPhotoListLogic) FindPhotoList(req *types.PhotoQueryReq) (resp *types.PageResp, err error) {
	in := &photorpc.FindPhotoListReq{
		AlbumId: req.AlbumId,
	}
	out, err := l.svcCtx.PhotoRpc.FindPhotoList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Photo, 0)
	for _, v := range out.List {
		m := ConvertPhotoTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = int64(len(list))
	resp.List = list
	return resp, nil
}

func ConvertPhotoTypes(req *photorpc.PhotoDetails) (out *types.Photo) {

	return &types.Photo{
		Id:       req.Id,
		PhotoUrl: req.PhotoSrc,
	}
}
