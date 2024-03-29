package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type OperationLogService struct {
	svcCtx *svc.ServiceContext
}

func NewOperationLogService(svcCtx *svc.ServiceContext) *OperationLogService {
	return &OperationLogService{
		svcCtx: svcCtx,
	}
}

// 创建OperationLog记录
func (s *OperationLogService) CreateOperationLog(reqCtx *request.Context, operationLog *entity.OperationLog) (data *entity.OperationLog, err error) {
	return s.svcCtx.OperationLogRepository.Create(reqCtx, operationLog)
}

// 更新OperationLog记录
func (s *OperationLogService) UpdateOperationLog(reqCtx *request.Context, operationLog *entity.OperationLog) (data *entity.OperationLog, err error) {
	return s.svcCtx.OperationLogRepository.Update(reqCtx, operationLog)
}

// 删除OperationLog记录
func (s *OperationLogService) DeleteOperationLog(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.OperationLogRepository.Delete(reqCtx, "id = ?", id)
}

// 查询OperationLog记录
func (s *OperationLogService) FindOperationLog(reqCtx *request.Context, id int) (data *entity.OperationLog, err error) {
	return s.svcCtx.OperationLogRepository.First(reqCtx, "id = ?", id)
}

// 批量删除OperationLog记录
func (s *OperationLogService) DeleteOperationLogByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.OperationLogRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取OperationLog记录
func (s *OperationLogService) FindOperationLogList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.OperationLog, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.OperationLogRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.OperationLogRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
