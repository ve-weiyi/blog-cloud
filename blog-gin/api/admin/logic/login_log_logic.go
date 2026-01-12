package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type LoginLogLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewLoginLogLogic(svcCtx *svctx.ServiceContext) *LoginLogLogic {
	return &LoginLogLogic{
		svcCtx: svcCtx,
	}
}

// 查询登录日志
func (s *LoginLogLogic) FindLoginLogList(reqCtx *request.Context, in *types.QueryLoginLogReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 删除登录日志
func (s *LoginLogLogic) DeletesLoginLog(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}
