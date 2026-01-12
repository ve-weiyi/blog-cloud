package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindVisitLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindVisitLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVisitLogListLogic {
	return &FindVisitLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询操作访问列表
func (l *FindVisitLogListLogic) FindVisitLogList(in *syslogrpc.FindVisitLogListReq) (*syslogrpc.FindVisitLogListResp, error) {
	page, size, sorts, conditions, params := convertVisitLogQuery(in)

	records, total, err := l.svcCtx.TVisitLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.VisitLog
	for _, v := range records {
		list = append(list, convertVisitLogOut(v))
	}

	return &syslogrpc.FindVisitLogListResp{
		List: list,
		Pagination: &syslogrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertVisitLogQuery(in *syslogrpc.FindVisitLogListReq) (page int, size int, sorts string, conditions string, params []any) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.UserId != "" {
		opts = append(opts, query.WithCondition("user_id = ?", in.UserId))
	}

	if in.TerminalId != "" {
		opts = append(opts, query.WithCondition("terminal_id = ?", in.TerminalId))
	}

	if in.PageName != "" {
		opts = append(opts, query.WithCondition("page like ?", "%"+in.PageName+"%"))
	}

	return query.NewQueryBuilder(opts...).Build()
}

func convertVisitLogOut(in *model.TVisitLog) (out *syslogrpc.VisitLog) {
	out = &syslogrpc.VisitLog{
		Id:         in.Id,
		UserId:     in.UserId,
		TerminalId: in.TerminalId,
		PageName:   in.PageName,
		CreatedAt:  in.CreatedAt.UnixMilli(),
		UpdatedAt:  in.UpdatedAt.UnixMilli(),
	}

	return out
}
