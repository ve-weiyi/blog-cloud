package banner

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
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

func (l *FindBannerListLogic) FindBannerList(req *types.BannerQueryReq) (resp *types.PageResp, err error) {
	in := &photorpc.FindBannerListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	out, err := l.svcCtx.PhotoRpc.FindBannerList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Banner, 0)
	for _, v := range out.List {
		list = append(list, ConvertBannerTypes(v))
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertBannerTypes(in *photorpc.BannerDetails) *types.Banner {
	return &types.Banner{
		Id:          in.Id,
		BannerName:  in.BannerName,
		BannerLabel: in.BannerLabel,
		BannerCover: in.BannerCover,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}
