package logic

import (
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/api/admin/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/infra/request"
	"github.com/ve-weiyi/ve-blog-golang/blog-gin/svctx"
)

type ApiLogic struct {
	svcCtx *svctx.ServiceContext
}

func NewApiLogic(svcCtx *svctx.ServiceContext) *ApiLogic {
	return &ApiLogic{
		svcCtx: svcCtx,
	}
}

// 创建api路由
func (s *ApiLogic) AddApi(reqCtx *request.Context, in *types.NewApiReq) (out *types.ApiBackVO, err error) {
	// todo

	return
}

// 清空接口列表
func (s *ApiLogic) CleanApiList(reqCtx *request.Context, in *types.EmptyReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 删除api路由
func (s *ApiLogic) DeletesApi(reqCtx *request.Context, in *types.IdsReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 分页获取api路由列表
func (s *ApiLogic) FindApiList(reqCtx *request.Context, in *types.QueryApiReq) (out *types.PageResp, err error) {
	// todo

	return
}

// 同步api列表
func (s *ApiLogic) SyncApiList(reqCtx *request.Context, in *types.SyncApiReq) (out *types.BatchResp, err error) {
	// todo

	return
}

// 更新api路由
func (s *ApiLogic) UpdateApi(reqCtx *request.Context, in *types.NewApiReq) (out *types.ApiBackVO, err error) {
	// todo

	return
}
