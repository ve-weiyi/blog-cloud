package account

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindAccountLoginHistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户登录历史
func NewFindAccountLoginHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindAccountLoginHistoryListLogic {
	return &FindAccountLoginHistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindAccountLoginHistoryListLogic) FindAccountLoginHistoryList(req *types.AccountQuery) (resp *types.PageResp, err error) {
	in := &accountrpc.FindLoginHistoryListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	out, err := l.svcCtx.AccountRpc.FindUserLoginHistoryList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}
	// 查询用户信息
	users, err := l.svcCtx.AccountRpc.FindUserList(l.ctx, &accountrpc.FindUserListReq{
		UserIds: uids,
	})
	if err != nil {
		return nil, err
	}

	usm := make(map[string]*accountrpc.User)
	for _, v := range users.List {
		usm[v.UserId] = v
	}

	var list []*types.AccountLoginHistory
	for _, v := range out.List {
		m := ConvertUserLoginHistoryTypes(v, usm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertUserLoginHistoryTypes(in *accountrpc.UserLoginHistory, usm map[string]*accountrpc.User) (out *types.AccountLoginHistory) {
	out = &types.AccountLoginHistory{
		Id:        in.Id,
		LoginType: in.LoginType,
		Agent:     in.Agent,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		LoginAt:   in.LoginAt,
		LogoutAt:  in.LogoutAt,
	}

	// 用户信息
	if in.UserId != "" {
		user, ok := usm[in.UserId]
		if ok && user != nil {
			out.Username = user.Username
			out.Nickname = user.Nickname
			out.Avatar = user.Avatar
		}
	}

	return
}
