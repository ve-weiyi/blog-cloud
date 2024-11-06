package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type OperationLogService struct {
	svcCtx *svctx.ServiceContext
}

func NewOperationLogService(svcCtx *svctx.ServiceContext) *OperationLogService {
	return &OperationLogService{
		svcCtx: svcCtx,
	}
}

// 批量删除操作记录
func (s *OperationLogService) BatchDeleteOperationLog(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除操作记录
func (s *OperationLogService) DeleteOperationLog(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 分页获取操作记录列表
func (s *OperationLogService) FindOperationLogList(reqCtx *request.Context, in *dto.OperationLogQuery) (out *dto.PageResp, err error) {
	// todo

	return
}
