package syslogrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOperationLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOperationLogListLogic {
	return &FindOperationLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取操作记录列表
func (l *FindOperationLogListLogic) FindOperationLogList(in *syslogrpc.FindOperationLogListReq) (*syslogrpc.FindOperationLogListResp, error) {
	page, size, sorts, conditions, params := convertOperationLogQuery(in)

	records, total, err := l.svcCtx.TOperationLogModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.OperationLog
	for _, v := range records {
		list = append(list, convertOperationLogOut(v))
	}

	return &syslogrpc.FindOperationLogListResp{
		List: list,
		Pagination: &syslogrpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertOperationLogQuery(in *syslogrpc.FindOperationLogListReq) (page int, size int, sorts string, conditions string, params []any) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.Keywords != "" {
		opts = append(opts, query.WithCondition("opt_desc = ?", "%"+in.Keywords+"%"))
	}

	return query.NewQueryBuilder(opts...).Build()
}

func convertOperationLogOut(in *model.TOperationLog) (out *syslogrpc.OperationLog) {
	out = &syslogrpc.OperationLog{
		Id:             in.Id,
		UserId:         in.UserId,
		TerminalId:     in.TerminalId,
		OptModule:      in.OptModule,
		OptDesc:        in.OptDesc,
		RequestUri:     in.RequestUri,
		RequestMethod:  in.RequestMethod,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}

	return out
}
