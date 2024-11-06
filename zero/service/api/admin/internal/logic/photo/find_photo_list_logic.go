package photo

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

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

func (l *FindPhotoListLogic) FindPhotoList(req *types.PhotoQuery) (resp *types.PageResp, err error) {
	in := &photorpc.FindPhotoListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
		AlbumId:  req.AlbumId,
	}

	out, err := l.svcCtx.PhotoRpc.FindPhotoList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.PhotoBackDTO
	for _, v := range out.List {
		m := ConvertPhotoTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
