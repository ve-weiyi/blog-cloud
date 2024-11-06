package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserApisLogic {
	return &FindUserApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户接口权限
func (l *FindUserApisLogic) FindUserApis(in *permissionrpc.UserIdReq) (*permissionrpc.FindApiListResp, error) {
	uid := in.UserId

	// 查用户
	// ua, err := l.svcCtx.TUserModel.First(l.ctx, "id = ?", uid)
	// if err != nil {
	//	return nil, err
	// }

	// 查用户角色
	urs, err := l.svcCtx.TUserRoleModel.FindALL(l.ctx, "user_id = ?", uid)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0)
	for _, v := range urs {
		ids = append(ids, v.RoleId)
	}

	// 查角色拥有的接口
	rs, err := l.svcCtx.TRoleApiModel.FindALL(l.ctx, "id in (?)", ids)
	if err != nil {
		return nil, err
	}

	var apiIds []int64
	for _, v := range rs {
		apiIds = append(apiIds, v.ApiId)
	}

	// 查接口信息
	apis, err := l.svcCtx.TApiModel.FindALL(l.ctx, "id in (?)", apiIds)
	if err != nil {
		return nil, err
	}

	var list []*permissionrpc.ApiDetails
	for _, v := range apis {
		list = append(list, convertApiOut(v))
	}

	out := &permissionrpc.FindApiListResp{}
	out.List = list
	return out, nil
}
