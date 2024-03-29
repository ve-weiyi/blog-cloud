package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
)

type UserAccountService struct {
	svcCtx *svc.ServiceContext
}

func NewUserAccountService(svcCtx *svc.ServiceContext) *UserAccountService {
	return &UserAccountService{
		svcCtx: svcCtx,
	}
}

// 创建UserAccount记录
func (s *UserAccountService) CreateUserAccount(reqCtx *request.Context, userAccount *entity.UserAccount) (data *entity.UserAccount, err error) {
	return s.svcCtx.UserAccountRepository.Create(reqCtx, userAccount)
}

// 更新UserAccount记录
func (s *UserAccountService) UpdateUserAccount(reqCtx *request.Context, userAccount *entity.UserAccount) (data *entity.UserAccount, err error) {
	return s.svcCtx.UserAccountRepository.Update(reqCtx, userAccount)
}

// 删除UserAccount记录
func (s *UserAccountService) DeleteUserAccount(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.UserAccountRepository.Delete(reqCtx, "id = ?", id)
}

// 查询UserAccount记录
func (s *UserAccountService) FindUserAccount(reqCtx *request.Context, id int) (data *entity.UserAccount, err error) {
	return s.svcCtx.UserAccountRepository.First(reqCtx, "id = ?", id)
}

// 批量删除UserAccount记录
func (s *UserAccountService) DeleteUserAccountByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.UserAccountRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取UserAccount记录
func (s *UserAccountService) FindUserAccountList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.UserAccount, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.UserAccountRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.UserAccountRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
