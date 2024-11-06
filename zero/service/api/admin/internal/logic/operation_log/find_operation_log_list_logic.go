package operation_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindOperationLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取操作记录列表
func NewFindOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOperationLogListLogic {
	return &FindOperationLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindOperationLogListLogic) FindOperationLogList(req *types.OperationLogQuery) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindOperationLogListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
	}

	out, err := l.svcCtx.SyslogRpc.FindOperationLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.OperationLogBackDTO
	for _, v := range out.List {
		m := ConvertOperationLogTypes(v)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertOperationLogTypes(in *syslogrpc.OperationLogDetails) (out *types.OperationLogBackDTO) {

	return &types.OperationLogBackDTO{
		Id:             in.Id,
		UserId:         in.UserId,
		Nickname:       in.Nickname,
		IpAddress:      in.IpAddress,
		IpSource:       in.IpSource,
		OptModule:      in.OptModule,
		OptDesc:        in.OptDesc,
		RequestUrl:     in.RequestUrl,
		RequestMethod:  in.RequestMethod,
		RequestHeader:  in.RequestHeader,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		ResponseStatus: in.ResponseStatus,
		Cost:           in.Cost,
		CreatedAt:      in.CreatedAt,
		UpdatedAt:      in.UpdatedAt,
	}
}
