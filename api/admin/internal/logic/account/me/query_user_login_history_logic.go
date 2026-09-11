package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type QueryUserLoginHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户登录历史
func NewQueryUserLoginHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryUserLoginHistoryLogic {
	return &QueryUserLoginHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryUserLoginHistoryLogic) QueryUserLoginHistory(req *types.QueryUserLoginHistoryReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.SyslogService.ListLoginLogs(l.ctx, &syslogservice.ListLoginLogsRequest{
		PageQuery: &syslogservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
	})
	if err != nil {
		return nil, err
	}

	var list []*types.LoginLogVO
	for _, v := range out.List {
		list = append(list, &types.LoginLogVO{
			Id:        v.Id,
			UserId:    v.UserId,
			DeviceId:  v.DeviceId,
			LoginType: v.LoginType,
			LoginAt:   v.CreatedAt,
			LogoutAt:  v.LogoutAt,
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
