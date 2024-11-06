package banner

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateBannerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新页面
func NewUpdateBannerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateBannerLogic {
	return &UpdateBannerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBannerLogic) UpdateBanner(req *types.BannerNewReq) (resp *types.BannerBackDTO, err error) {
	in := ConvertBannerPb(req)
	out, err := l.svcCtx.PhotoRpc.UpdateBanner(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = ConvertBannerTypes(out)
	return resp, nil
}
