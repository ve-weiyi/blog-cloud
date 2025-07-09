package permissionrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/permissionrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserMenusLogic {
	return &FindUserMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户菜单权限
func (l *FindUserMenusLogic) FindUserMenus(in *permissionrpc.UserIdReq) (*permissionrpc.FindMenuListResp, error) {
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

	// 查角色拥有的菜单
	rs, err := l.svcCtx.TRoleMenuModel.FindALL(l.ctx, "role_id in (?)", ids)
	if err != nil {
		return nil, err
	}

	var mids []int64
	for _, v := range rs {
		mids = append(mids, v.MenuId)
	}

	// 查菜单信息
	list, err := l.svcCtx.TMenuModel.FindALL(l.ctx, "id in (?)", mids)
	if err != nil {
		return nil, err
	}

	var root permissionrpc.MenuDetailsResp
	root.Children = appendMenuChildren(&root, list)

	out := &permissionrpc.FindMenuListResp{}
	out.List = root.Children

	return out, nil
}
