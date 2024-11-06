package banner

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/photorpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建页面
func NewAddBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddBannerLogic {
	return &AddBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBannerLogic) AddBanner(req *types.BannerNewReq) (resp *types.BannerBackDTO, err error) {
	in := ConvertBannerPb(req)
	out, err := l.svcCtx.PhotoRpc.AddBanner(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertBannerTypes(out)
	return resp, nil
}

func ConvertBannerPb(in *types.BannerNewReq) (out *photorpc.BannerNewReq) {
	out = &photorpc.BannerNewReq{
		Id:          in.Id,
		BannerName:  in.BannerName,
		BannerLabel: in.BannerLabel,
		BannerCover: in.BannerCover,
	}

	return
}

func ConvertBannerTypes(in *photorpc.BannerDetails) (out *types.BannerBackDTO) {
	out = &types.BannerBackDTO{
		Id:          in.Id,
		BannerName:  in.BannerName,
		BannerLabel: in.BannerLabel,
		BannerCover: in.BannerCover,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}

	return
}
