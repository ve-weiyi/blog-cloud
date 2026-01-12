package visit_log

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindVisitLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取操作记录列表
func NewFindVisitLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVisitLogListLogic {
	return &FindVisitLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindVisitLogListLogic) FindVisitLogList(req *types.QueryVisitLogReq) (resp *types.PageResp, err error) {
	in := &syslogrpc.FindVisitLogListReq{
		Paginate: &syslogrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		UserId:     req.UserId,
		TerminalId: req.TerminalId,
		PageName:   req.PageName,
	}

	out, err := l.svcCtx.SyslogRpc.FindVisitLogList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	usm, err := apiutils.BatchQuery(out.List,
		func(v *syslogrpc.VisitLog) string {
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
		func(v *syslogrpc.VisitLog) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	var list []*types.VisitLogBackVO
	for _, v := range out.List {
		m := &types.VisitLogBackVO{
			Id:         v.Id,
			UserId:     v.UserId,
			TerminalId: v.TerminalId,
			PageName:   v.PageName,
			CreatedAt:  v.CreatedAt,
			UpdatedAt:  v.UpdatedAt,
			UserInfo:   usm[v.UserId],
			ClientInfo: vsm[v.TerminalId],
		}
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
