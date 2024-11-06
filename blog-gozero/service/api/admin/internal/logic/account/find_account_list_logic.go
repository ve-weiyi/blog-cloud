package account

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAccountListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户列表
func NewFindAccountListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAccountListLogic {
	return &FindAccountListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAccountListLogic) FindAccountList(req *types.AccountQuery) (resp *types.PageResp, err error) {
	in := &accountrpc.FindUserListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    "",
		Phone:    "",
		Status:   0,
		UserIds:  nil,
	}

	out, err := l.svcCtx.AccountRpc.FindUserInfoList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserInfoResp
	for _, v := range out.List {
		m := ConvertUserInfoTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertUserInfoTypes(in *accountrpc.UserInfoResp) *types.UserInfoResp {
	roles := make([]*types.UserRoleLabel, 0)
	for _, v := range in.Roles {
		m := &types.UserRoleLabel{
			RoleId:      v.RoleId,
			RoleName:    v.RoleName,
			RoleComment: v.RoleComment,
		}

		roles = append(roles, m)
	}

	var info types.UserInfoExt
	jsonconv.JsonToAny(in.Info, &info)

	out := &types.UserInfoResp{
		UserId:      in.UserId,
		Username:    in.Username,
		Nickname:    in.Nickname,
		Avatar:      in.Avatar,
		Email:       in.Email,
		Phone:       in.Phone,
		Status:      in.Status,
		LoginType:   in.LoginType,
		IpAddress:   in.IpAddress,
		IpSource:    in.IpSource,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
		Roles:       roles,
		UserInfoExt: info,
	}

	return out
}
