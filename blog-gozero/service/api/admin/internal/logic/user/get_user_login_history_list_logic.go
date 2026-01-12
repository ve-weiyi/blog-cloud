package user

import (
	"context"

	"github.com/spf13/cast"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/pkg/infra/biz/bizheader"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
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

func (l *GetUserLoginHistoryListLogic) GetUserLoginHistoryList(req *types.QueryUserLoginHistoryReq) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindLoginLogListReq{
		Paginate: &syslogrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		UserId: cast.ToString(l.ctx.Value(bizheader.HeaderUid)),
	}

	out, err := l.svcCtx.SyslogRpc.FindLoginLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.UserLoginHistory
	for _, v := range out.List {
		list = append(list, &types.UserLoginHistory{
			Id:         v.Id,
			UserId:     v.UserId,
			TerminalId: v.TerminalId,
			LoginType:  v.LoginType,
			AppName:    v.AppName,
			LoginAt:    v.LoginAt,
			LogoutAt:   v.LogoutAt,
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
