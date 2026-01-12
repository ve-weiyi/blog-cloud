package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *FindPhotoListLogic) FindPhotoList(req *types.QueryPhotoReq) (resp *types.PageResp, err error) {
	in := &resourcerpc.FindPhotoListReq{
		Paginate: &resourcerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		AlbumId:  req.AlbumId,
		IsDelete: req.IsDelete,
	}

	out, err := l.svcCtx.ResourceRpc.FindPhotoList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.PhotoBackVO
	for _, v := range out.List {
		list = append(list, &types.PhotoBackVO{
			Id:        v.Id,
			AlbumId:   v.AlbumId,
			PhotoName: v.PhotoName,
			PhotoDesc: v.PhotoDesc,
			PhotoSrc:  v.PhotoSrc,
			IsDelete:  v.IsDelete,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
