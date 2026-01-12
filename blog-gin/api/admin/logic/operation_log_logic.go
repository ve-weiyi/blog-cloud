package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type OperationLogLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewOperationLogLogic(svcCtx *svctx.ServiceContext) *OperationLogLogic {
	return &OperationLogLogic{
		svcCtx: svcCtx,
	}
}

// 删除操作记录
func (s *OperationLogLogic) DeletesOperationLog(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取操作记录列表
func (s *OperationLogLogic) FindOperationLogList(reqCtx *request.Context, in *types.QueryOperationLogReq) (out *types.PageResp, err error) {
	// todo

	return
}
