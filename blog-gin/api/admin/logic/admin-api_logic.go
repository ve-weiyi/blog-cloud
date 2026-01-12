package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type AdminApiLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewAdminApiLogic(svcCtx *svctx.ServiceContext) *AdminApiLogic {
	return &AdminApiLogic{
		svcCtx: svcCtx,
	}
}

// ping
func (s *AdminApiLogic) Ping(reqCtx *request.Context, in *types.PingReq) (out *types.PingResp, err error) {
	// todo

	return
}
