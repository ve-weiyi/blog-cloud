package service

import (
	"github.com/ve-weiyi/ve-blog-golang/server/api/admin/dto"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/base/request"
	"github.com/ve-weiyi/ve-blog-golang/server/svctx"
)

type ApiService struct {
	svcCtx *svctx.ServiceContext
}

func NewApiService(svcCtx *svctx.ServiceContext) *ApiService {
	return &ApiService{
		svcCtx: svcCtx,
	}
}

// 分页获取api路由列表
func (s *ApiService) FindApiList(reqCtx *request.Context, in *dto.ApiQuery) (out *dto.PageResp, err error) {
	// todo

	return
}

// 创建api路由
func (s *ApiService) AddApi(reqCtx *request.Context, in *dto.ApiNewReq) (out *dto.ApiBackDTO, err error) {
	// todo

	return
}

// 批量删除api路由
func (s *ApiService) BatchDeleteApi(reqCtx *request.Context, in *dto.IdsReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 清空接口列表
func (s *ApiService) CleanApiList(reqCtx *request.Context, in *dto.EmptyReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 删除api路由
func (s *ApiService) DeleteApi(reqCtx *request.Context, in *dto.IdReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 同步api列表
func (s *ApiService) SyncApiList(reqCtx *request.Context, in *dto.SyncApiReq) (out *dto.BatchResp, err error) {
	// todo

	return
}

// 更新api路由
func (s *ApiService) UpdateApi(reqCtx *request.Context, in *dto.ApiNewReq) (out *dto.ApiBackDTO, err error) {
	// todo

	return
}
