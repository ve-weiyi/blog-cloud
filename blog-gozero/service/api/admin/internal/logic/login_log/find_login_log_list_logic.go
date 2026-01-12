package login_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindLoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户登录历史
func NewFindLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindLoginLogListLogic {
	return &FindLoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindLoginLogListLogic) FindLoginLogList(req *types.QueryLoginLogReq) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindLoginLogListReq{
		Paginate: &syslogrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		UserId: req.UserId,
	}

	out, err := l.svcCtx.SyslogRpc.FindLoginLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	usm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.LoginLog) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}
	// 查询访客信息
	vsm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.LoginLog) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}
	var list []*types.LoginLogBackVO
	for _, v := range out.List {
		list = append(list, &types.LoginLogBackVO{
			Id:         v.Id,
			UserId:     v.UserId,
			TerminalId: v.TerminalId,
			LoginType:  v.LoginType,
			AppName:    v.AppName,
			LoginAt:    v.LoginAt,
			LogoutAt:   v.LogoutAt,
			UserInfo:   usm[v.UserId],
			ClientInfo: vsm[v.TerminalId],
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
