package visit_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type QueryVisitLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取访问日志列表
func NewQueryVisitLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryVisitLogListLogic {
	return &QueryVisitLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryVisitLogListLogic) QueryVisitLogList(req *types.QueryVisitLogListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.SyslogService.ListVisitLogs(l.ctx, &syslogservice.ListVisitLogsRequest{
		PageQuery: &syslogservice.PageQuery{Page: req.Page, PageSize: req.PageSize, Sorts: req.Sorts},
		UserId:    req.UserId,
		DeviceId:  req.DeviceId,
		PageName:  req.PageName,
	})
	if err != nil {
		return nil, err
	}

	var uids, tids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
		tids = append(tids, v.DeviceId)
	}

	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	vsm, err := apiutils.GetGuests(l.ctx, l.svcCtx, tids)
	if err != nil {
		return nil, err
	}

	var list []*types.VisitLogVO
	for _, v := range out.List {
		list = append(list, &types.VisitLogVO{
			Id:        v.Id,
			UserId:    v.UserId,
			DeviceId:  v.DeviceId,
			PageName:  v.PageName,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			UserInfo:  usm[v.UserId],
			GuestInfo: vsm[v.DeviceId],
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
