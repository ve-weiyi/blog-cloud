package visitor

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/accountrpc"
)

type FindVisitorListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取游客列表
func NewFindVisitorListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindVisitorListLogic {
	return &FindVisitorListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindVisitorListLogic) FindVisitorList(req *types.QueryVisitorReq) (resp *types.PageResp, err error) {
	in := &accountrpc.FindVisitorListReq{
		Paginate: &accountrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TerminalId: req.TerminalId,
		IpSource:   req.IpSource,
	}

	out, err := l.svcCtx.AccountRpc.FindVisitorList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var list []*types.VisitorBackVO
	for _, item := range out.List {
		list = append(list, &types.VisitorBackVO{
			Id:         item.Id,
			TerminalId: item.TerminalId,
			Os:         item.Os,
			Browser:    item.Browser,
			IpAddress:  item.IpAddress,
			IpSource:   item.IpSource,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
