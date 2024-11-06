package syslogrpclogic

import (
	"context"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/syslogrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

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

	result, err := l.svcCtx.TOperationLogModel.FindList(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	count, err := l.svcCtx.TOperationLogModel.FindCount(l.ctx, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*syslogrpc.OperationLogDetails
	for _, v := range result {
		list = append(list, convertOperationLogOut(v))
	}

	return &syslogrpc.FindOperationLogListResp{
		List:  list,
		Total: count,
	}, nil
}

func convertOperationLogQuery(in *syslogrpc.FindOperationLogListReq) (page int, size int, sorts string, conditions string, params []any) {
	page = int(in.Page)
	size = int(in.PageSize)
	sorts = strings.Join(in.Sorts, ",")
	if sorts == "" {
		sorts = "id desc"
	}

	return
}
