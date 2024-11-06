package banner

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindBannerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取页面列表
func NewFindBannerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindBannerListLogic {
	return &FindBannerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindBannerListLogic) FindBannerList(req *types.BannerQuery) (resp *types.PageResp, err error) {
	in := &photorpc.FindBannerListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
	}

	out, err := l.svcCtx.PhotoRpc.FindBannerList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.BannerBackDTO
	for _, v := range out.List {
		m := ConvertBannerTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
