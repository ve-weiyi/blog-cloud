package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/global/constant"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/configrpc"
	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAboutMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新关于我的信息
func NewUpdateAboutMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAboutMeLogic {
	return &UpdateAboutMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAboutMeLogic) UpdateAboutMe(req *types.AboutMeVO) (resp *types.EmptyResp, err error) {
	in := &configrpc.SaveConfigReq{
		ConfigKey:   constant.ConfigKeyAboutMe,
		ConfigValue: jsonconv.AnyToJsonNE(req),
	}

	_, err = l.svcCtx.ConfigRpc.SaveConfig(l.ctx, in)
	if err != nil {
		return nil, err
	}

	return
}
