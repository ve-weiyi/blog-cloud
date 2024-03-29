package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type ApiService struct {
	svcCtx *svc.ServiceContext
}

func NewApiService(svcCtx *svc.ServiceContext) *ApiService {
	return &ApiService{
		svcCtx: svcCtx,
	}
}

// 创建Api记录
func (s *ApiService) CreateApi(reqCtx *request.Context, api *entity.Api) (data *entity.Api, err error) {
	return s.svcCtx.ApiRepository.Create(reqCtx, api)
}

// 更新Api记录
func (s *ApiService) UpdateApi(reqCtx *request.Context, api *entity.Api) (data *entity.Api, err error) {
	return s.svcCtx.ApiRepository.Update(reqCtx, api)
}

// 删除Api记录
func (s *ApiService) DeleteApi(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.ApiRepository.Delete(reqCtx, "id = ?", id)
}

// 查询Api记录
func (s *ApiService) FindApi(reqCtx *request.Context, id int) (data *entity.Api, err error) {
	return s.svcCtx.ApiRepository.First(reqCtx, "id = ?", id)
}

// 批量删除Api记录
func (s *ApiService) DeleteApiByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.ApiRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取Api记录
func (s *ApiService) FindApiList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.Api, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.ApiRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.ApiRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
