package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *accountrpc.UserIdReq) (*accountrpc.UserInfoResp, error) {
	uid := in.UserId

	ui, err := l.svcCtx.TUserModel.First(l.ctx, "id = ?", uid)
	if err != nil {
		return nil, err
	}

	// 查找用户角色
	rList, err := getUserRoles(l.ctx, l.svcCtx, uid)
	if err != nil {
		return nil, err
	}

	return convertUserInfoOut(ui, rList), nil
}

func getUserRoles(ctx context.Context, svcCtx *svc.ServiceContext, uid string) (list []*model.TRole, err error) {
	// 查找用户角色
	urList, err := svcCtx.TUserRoleModel.FindALL(ctx, "user_id in (?)", uid)
	if err != nil {
		return nil, err
	}

	var ursMap = make(map[string][]int64)
	var roleIds []int64
	for _, item := range urList {
		roleIds = append(roleIds, item.RoleId)
		ursMap[item.UserId] = append(ursMap[item.UserId], item.RoleId)
	}

	// 查找角色信息
	rList, err := svcCtx.TRoleModel.FindALL(ctx, "id in (?)", roleIds)
	if err != nil {
		return nil, err
	}

	return rList, nil
}

func convertUserInfoOut(in *model.TUser, roles []*model.TRole) (out *accountrpc.UserInfoResp) {
	var list []*accountrpc.UserRoleLabel
	for _, role := range roles {
		m := &accountrpc.UserRoleLabel{
			RoleId:      role.Id,
			RoleName:    role.RoleName,
			RoleComment: role.RoleComment,
		}

		list = append(list, m)
	}

	out = &accountrpc.UserInfoResp{
		UserId:    in.UserId,
		Username:  in.Username,
		Nickname:  in.Nickname,
		Avatar:    in.Avatar,
		Email:     in.Email,
		Phone:     in.Phone,
		Info:      in.Info,
		Status:    in.Status,
		LoginType: in.LoginType,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		CreatedAt: in.CreatedAt.Unix(),
		UpdatedAt: in.UpdatedAt.Unix(),
		Roles:     list,
	}

	return out
}
