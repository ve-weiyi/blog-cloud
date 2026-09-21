package syslogservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/infra/queryx"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type ListOperationLogsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOperationLogsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperationLogsLogic {
	return &ListOperationLogsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询操作日志
func (l *ListOperationLogsLogic) ListOperationLogs(in *syslogrpc.ListOperationLogsRequest) (*syslogrpc.ListOperationLogsResponse, error) {
	var opts []queryx.Option
	if in.ListQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.ListQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.ListQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.ListQuery.Sorts...))
	}
	if in.UserId != nil {
		opts = append(opts, queryx.WithCondition("user_id = ?", *in.UserId))
	}
	if in.Module != nil {
		opts = append(opts, queryx.WithCondition("module = ?", *in.Module))
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
	records, total, err := l.svcCtx.TOperationLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.OperationLog
	for _, v := range records {
		list = append(list, convertOperationLogOut(v))
	}

	return &syslogrpc.ListOperationLogsResponse{
		ListResult: &syslogrpc.ListResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
