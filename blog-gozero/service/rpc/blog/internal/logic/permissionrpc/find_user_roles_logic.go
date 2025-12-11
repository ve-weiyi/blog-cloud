package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserRolesLogic {
	return &FindUserRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户角色信息
func (l *FindUserRolesLogic) FindUserRoles(in *permissionrpc.UserIdReq) (*permissionrpc.FindRoleListResp, error) {
	uid := in.UserId

	// 查用户
	// ua, err := l.svcCtx.TUserModel.FindOneByUserId(l.ctx, uid)
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

	// 查角色
	rs, err := l.svcCtx.TRoleModel.FindALL(l.ctx, "id in (?)", ids)
	if err != nil {
		return nil, err
	}

	var list []*permissionrpc.RoleDetailsResp
	for _, v := range rs {
		list = append(list, convertRoleOut(v))
	}

	out := &permissionrpc.FindRoleListResp{}
	out.List = list
	return out, nil
}
