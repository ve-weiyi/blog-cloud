package logic

import (
	"fmt"

	"github.com/ve-weiyi/ve-blog-golang/server/api/model/entity"
	"github.com/ve-weiyi/ve-blog-golang/server/api/model/request"
	"github.com/ve-weiyi/ve-blog-golang/server/api/service/svc"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/mail"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

type RemarkService struct {
	svcCtx *svc.ServiceContext
}

func NewRemarkService(svcCtx *svc.ServiceContext) *RemarkService {
	return &RemarkService{
		svcCtx: svcCtx,
	}
}

// 创建Remark记录
func (s *RemarkService) CreateRemark(reqCtx *request.Context, remark *entity.Remark) (data *entity.Remark, err error) {
	remark.IpAddress = reqCtx.IpAddress
	remark.IpSource = reqCtx.GetIpSource()
	msg := &mail.EmailMessage{
		To:      []string{"791422171@qq.com"},
		Subject: "【blog】新增留言信息",
		Content: fmt.Sprintf("%v remark: %v", reqCtx.Username, remark.MessageContent),
		Type:    0,
	}

	err = s.svcCtx.EmailPublisher.PublishMessage([]byte(jsonconv.ObjectToJson(msg)))
	if err != nil {
		s.svcCtx.Log.Info("PublishMessage:", err)
	}

	return s.svcCtx.RemarkRepository.Create(reqCtx, remark)
}

// 更新Remark记录
func (s *RemarkService) UpdateRemark(reqCtx *request.Context, remark *entity.Remark) (data *entity.Remark, err error) {
	return s.svcCtx.RemarkRepository.Update(reqCtx, remark)
}

// 删除Remark记录
func (s *RemarkService) DeleteRemark(reqCtx *request.Context, id int) (rows int64, err error) {
	return s.svcCtx.RemarkRepository.Delete(reqCtx, "id = ?", id)
}

// 查询Remark记录
func (s *RemarkService) FindRemark(reqCtx *request.Context, id int) (data *entity.Remark, err error) {
	return s.svcCtx.RemarkRepository.First(reqCtx, "id = ?", id)
}

// 批量删除Remark记录
func (s *RemarkService) DeleteRemarkByIds(reqCtx *request.Context, ids []int) (rows int64, err error) {
	return s.svcCtx.RemarkRepository.Delete(reqCtx, "id in (?)", ids)
}

// 分页获取Remark记录
func (s *RemarkService) FindRemarkList(reqCtx *request.Context, page *request.PageQuery) (list []*entity.Remark, total int64, err error) {
	cond, args := page.ConditionClause()
	order := page.OrderClause()

	list, err = s.svcCtx.RemarkRepository.FindList(reqCtx, page.Page, page.PageSize, order, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	total, err = s.svcCtx.RemarkRepository.Count(reqCtx, cond, args...)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
