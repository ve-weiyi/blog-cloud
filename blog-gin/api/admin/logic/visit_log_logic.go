package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/common/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type VisitLogLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewVisitLogLogic(svcCtx *svctx.ServiceContext) *VisitLogLogic {
	return &VisitLogLogic{
		svcCtx: svcCtx,
	}
}

// 删除操作记录
func (s *VisitLogLogic) DeletesVisitLog(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取操作记录列表
func (s *VisitLogLogic) FindVisitLogList(reqCtx *request.Context, in *types.VisitLogQuery) (out *types.PageResp, err error) {
	// todo

	return
}
