package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLoginHistoryListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户登录历史
func NewGetUserLoginHistoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLoginHistoryListLogic {
	return &GetUserLoginHistoryListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserLoginHistoryListLogic) GetUserLoginHistoryList(req *types.UserLoginHistoryQuery) (resp *types.PageResp, err error) {
	in := &accountrpc.FindLoginHistoryListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		UserId:   cast.ToString(l.ctx.Value("uid")),
	}

	out, err := l.svcCtx.AccountRpc.FindUserLoginHistoryList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserLoginHistory
	for _, v := range out.List {
		m := ConvertUserLoginHistoryTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertUserLoginHistoryTypes(in *accountrpc.UserLoginHistory) *types.UserLoginHistory {
	return &types.UserLoginHistory{
		Id:        in.Id,
		LoginType: in.LoginType,
		Agent:     in.Agent,
		IpAddress: in.IpAddress,
		IpSource:  in.IpSource,
		LoginAt:   in.LoginAt,
		LogoutAt:  in.LogoutAt,
	}
}
