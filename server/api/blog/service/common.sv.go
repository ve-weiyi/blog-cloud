package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/blog/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type CommonService struct {
	svcCtx *svctx.ServiceContext
}

func NewCommonService(svcCtx *svctx.ServiceContext) *CommonService {
	return &CommonService{
		svcCtx: svcCtx,
	}
}

// ping
func (s *CommonService) Ping(reqCtx *request.Context, in *dto.PingReq) (out *dto.PingResp, err error) {
	// todo

	return
}
