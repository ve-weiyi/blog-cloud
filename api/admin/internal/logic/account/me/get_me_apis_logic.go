package me

import (
	"context"

	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/infra/biz/bizheader"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/accessservice"
)

type GetMeApisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户接口权限
func NewGetMeApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMeApisLogic {
	return &GetMeApisLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMeApisLogic) GetMeApis(req *types.EmptyReq) (resp *types.GetMeApisResp, err error) {
	userId := cast.ToString(l.ctx.Value(bizheader.HeaderUid))

	out, err := l.svcCtx.AccessService.GetUserApis(l.ctx, &accessservice.GetUserApisRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.UserApi
	for _, v := range out.List {
		list = append(list, convertUserApi(v))
	}

	return &types.GetMeApisResp{List: list}, nil
}

func convertUserApi(in *accessservice.Api) *types.UserApi {
	children := make([]*types.UserApi, 0)
	for _, v := range in.Children {
		children = append(children, convertUserApi(v))
	}

	return &types.UserApi{
		Id:        in.Id,
		ParentId:  in.ParentId,
		Name:      in.Name,
		Path:      in.Path,
		Method:    in.Method,
		CreatedAt: in.CreatedAt,
		UpdatedAt: in.UpdatedAt,
		Children:  children,
	}
}
