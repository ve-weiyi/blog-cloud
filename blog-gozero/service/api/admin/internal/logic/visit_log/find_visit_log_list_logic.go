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

func (l *FindVisitLogListLogic) FindVisitLogList(req *types.VisitLogQuery) (resp *types.PageResp, err error) {
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

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}

	// 查询用户信息
	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	var list []*types.VisitLogBackVO
	for _, v := range out.List {
		m := ConvertVisitLogTypes(v, usm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}

func ConvertVisitLogTypes(in *syslogrpc.VisitLogDetailsResp, usm map[string]*types.UserInfoVO) (out *types.VisitLogBackVO) {

	out = &types.VisitLogBackVO{
		Id:         in.Id,
		UserId:     in.UserId,
		TerminalId: in.TerminalId,
		PageName:   in.PageName,
		IpAddress:  in.IpAddress,
		IpSource:   in.IpSource,
		Os:         in.Os,
		Browser:    in.Browser,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}

	// 用户信息
	if in.UserId != "" {
		user, ok := usm[in.UserId]
		if ok && user != nil {
			out.User = user
		}
	}

	return out
}
