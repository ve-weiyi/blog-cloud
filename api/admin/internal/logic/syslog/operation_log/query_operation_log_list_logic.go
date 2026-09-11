package operation_log

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/admin/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/syslogservice"
)

type QueryOperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取操作日志列表
func NewQueryOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryOperationLogListLogic {
	return &QueryOperationLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryOperationLogListLogic) QueryOperationLogList(req *types.QueryOperationLogListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.SyslogService.ListOperationLogs(l.ctx, &syslogservice.ListOperationLogsRequest{
		PageQuery: &syslogservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
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

	var list []*types.OperationLogVO
	for _, v := range out.List {
		list = append(list, &types.OperationLogVO{
			Id:             v.Id,
			UserId:         v.UserId,
			DeviceId:       v.DeviceId,
			Module:         v.Module,
			Description:    v.Description,
			RequestUri:     v.RequestUri,
			RequestMethod:  v.RequestMethod,
			RequestData:    v.RequestData,
			ResponseData:   v.ResponseData,
			ResponseStatus: v.ResponseStatus,
			Cost:           v.Cost,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			UserInfo:       usm[v.UserId],
			GuestInfo:      vsm[v.DeviceId],
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
