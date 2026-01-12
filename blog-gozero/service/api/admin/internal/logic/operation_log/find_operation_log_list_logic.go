package operation_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

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

func (l *FindOperationLogListLogic) FindOperationLogList(req *types.QueryOperationLogReq) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindOperationLogListReq{
		Paginate: &syslogrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}

	out, err := l.svcCtx.SyslogRpc.FindOperationLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	usm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.OperationLogDetailsResp) string {
			return v.UserId
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	// 查询访客信息
	vsm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.OperationLogDetailsResp) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	var list []*types.OperationLogBackVO
	for _, v := range out.List {
		list = append(list, &types.OperationLogBackVO{
			Id:             v.Id,
			UserId:         v.UserId,
			TerminalId:     v.TerminalId,
			OptModule:      v.OptModule,
			OptDesc:        v.OptDesc,
			RequestUri:     v.RequestUri,
			RequestMethod:  v.RequestMethod,
			RequestData:    v.RequestData,
			ResponseData:   v.ResponseData,
			ResponseStatus: v.ResponseStatus,
			Cost:           v.Cost,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			UserInfo:       usm[v.UserId],
			ClientInfo:     vsm[v.TerminalId],
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
