package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/kit/infra/restx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/permissionrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserApisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户接口权限
func NewGetUserApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserApisLogic {
	return &GetUserApisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserApisLogic) GetUserApis(req *types.EmptyReq) (resp *types.UserApisResp, err error) {
	in := &permissionrpc.UserIdReq{
		UserId: cast.ToString(l.ctx.Value(restx.HeaderUid)),
	}

	out, err := l.svcCtx.PermissionRpc.FindUserApis(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserApi
	for _, v := range out.List {
		m := convertUserApi(v)
		list = append(list, m)
	}

	resp = &types.UserApisResp{}
	resp.List = list
	return
}

func convertUserApi(req *permissionrpc.ApiDetailsResp) (out *types.UserApi) {
	children := make([]*types.UserApi, 0)
	for _, v := range req.Children {
		m := convertUserApi(v)
		children = append(children, m)
	}

	out = &types.UserApi{
		Id:        req.Id,
		ParentId:  req.ParentId,
		Name:      req.Name,
		Path:      req.Path,
		Method:    req.Method,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.UpdatedAt,
		Children:  children,
	}

	return out
}
