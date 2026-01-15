package accountrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindVisitorListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindVisitorListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVisitorListLogic {
	return &FindVisitorListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询游客信息
func (l *FindVisitorListLogic) FindVisitorList(in *accountrpc.FindVisitorListReq) (*accountrpc.FindVisitorListResp, error) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.TerminalId != "" {
		opts = append(opts, query.WithCondition("terminal_id = ?", in.TerminalId))
	}

	if len(in.TerminalIds) != 0 {
		opts = append(opts, query.WithCondition("terminal_id in (?)", in.TerminalIds))
	}

	if in.IpSource != "" {
		opts = append(opts, query.WithCondition("ip_source like ?", "%"+in.IpSource+"%"))
	}

	page, size, sorts, conditions, params := query.NewQueryBuilder(opts...).Build()

	records, total, err := l.svcCtx.TVisitorModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*accountrpc.VisitorInfo
	for _, vs := range records {
		list = append(list, &accountrpc.VisitorInfo{
			Id:         vs.Id,
			TerminalId: vs.TerminalId,
			Os:         vs.Os,
			Browser:    vs.Browser,
			IpAddress:  vs.IpAddress,
			IpSource:   vs.IpSource,
			CreatedAt:  vs.CreatedAt.UnixMilli(),
			UpdatedAt:  vs.UpdatedAt.UnixMilli(),
		})
	}

	resp := &accountrpc.FindVisitorListResp{}
	resp.Pagination = &accountrpc.PageResp{
		Page:     int64(page),
		PageSize: int64(size),
		Total:    total,
	}
	resp.List = list

	return resp, nil
}
