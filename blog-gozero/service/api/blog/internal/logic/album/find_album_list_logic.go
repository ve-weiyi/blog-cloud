package album

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/resourcerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
)

type FindAlbumListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取相册列表
func NewFindAlbumListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAlbumListLogic {
	return &FindAlbumListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAlbumListLogic) FindAlbumList(req *types.QueryAlbumReq) (resp *types.PageResp, err error) {
	in := &resourcerpc.FindAlbumListReq{
		Paginate: &resourcerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}
	out, err := l.svcCtx.ResourceRpc.FindAlbumList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Album, 0)
	for _, v := range out.List {
		list = append(list, &types.Album{
			Id:         v.Id,
			AlbumName:  v.AlbumName,
			AlbumDesc:  v.AlbumDesc,
			AlbumCover: v.AlbumCover,
		})
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "相册",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
