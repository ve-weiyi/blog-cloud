package website

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/configrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAboutMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取关于我的信息
func NewGetAboutMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAboutMeLogic {
	return &GetAboutMeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAboutMeLogic) GetAboutMe(req *types.EmptyReq) (resp *types.AboutMe, err error) {
	in := &configrpc.FindConfigReq{
		ConfigKey: "about_me",
	}

	out, err := l.svcCtx.ConfigRpc.FindConfig(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.AboutMe{}
	jsonconv.JsonToAny(out.ConfigValue, &resp)
	return resp, nil
}
