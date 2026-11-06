package me

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/syslogservice"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
)

type QueryMeLoginLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询用户登录历史
func NewQueryMeLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryMeLoginLogListLogic {
	return &QueryMeLoginLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryMeLoginLogListLogic) QueryMeLoginLogList(req *types.QueryMeLoginLogListReq) (resp *types.ListResult, err error) {
	uid, err := metax.GetApiUserIdFromCtx(l.ctx)
	if err != nil || uid == "" {
		return nil, bizerr.NewBizError(bizcode.CodeUnauthenticated, "缺少用户身份")
	}

	out, err := l.svcCtx.SyslogService.ListLoginLogs(l.ctx, &syslogservice.ListLoginLogsRequest{
		ListQuery: &syslogservice.ListQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		UserId:    &uid,
	})
	if err != nil {
		return nil, err
	}

	var list []*types.LoginLogVO
	for _, v := range out.List {
		list = append(list, &types.LoginLogVO{
			Id:         v.Id,
			UserId:     v.UserId,
			DeviceId:   v.DeviceId,
			EventType:  v.EventType,
			LoginType:  v.LoginType,
			LogoutType: v.LogoutType,
			LoginAt:    v.CreatedAt,
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
