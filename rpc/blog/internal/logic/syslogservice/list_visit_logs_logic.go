package syslogservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListVisitLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListVisitLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListVisitLogsLogic {
	return &ListVisitLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询访问日志
func (l *ListVisitLogsLogic) ListVisitLogs(in *syslogrpc.ListVisitLogsRequest) (*syslogrpc.ListVisitLogsResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}
	if in.UserId != nil {
		opts = append(opts, queryx.WithCondition("user_id = ?", *in.UserId))
	}
	if in.DeviceId != nil {
		opts = append(opts, queryx.WithCondition("device_id = ?", *in.DeviceId))
	}
	if in.PageName != nil {
		opts = append(opts, queryx.WithCondition("page_name = ?", *in.PageName))
	}
	if in.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *in.StartDate); err == nil {
			opts = append(opts, queryx.WithCondition("created_at >= ?", t))
		}
	}
	if in.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *in.EndDate); err == nil {
			opts = append(opts, queryx.WithCondition("created_at <= ?", t.Add(24*time.Hour)))
		}
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TVisitLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.VisitLog
	for _, v := range records {
		list = append(list, convertVisitLogOut(v))
	}

	return &syslogrpc.ListVisitLogsResponse{
		PageResult: &syslogrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
