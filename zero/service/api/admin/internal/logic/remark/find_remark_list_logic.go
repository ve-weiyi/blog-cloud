package remark

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/messagerpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindRemarkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取留言列表
func NewFindRemarkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindRemarkListLogic {
	return &FindRemarkListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindRemarkListLogic) FindRemarkList(req *types.RemarkQuery) (resp *types.PageResp, err error) {
	in := &messagerpc.FindRemarkListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		Sorts:    req.Sorts,
	}

	out, err := l.svcCtx.MessageRpc.FindRemarkList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
	}

	// 查询用户信息
	users, err := l.svcCtx.AccountRpc.FindUserList(l.ctx, &accountrpc.FindUserListReq{
		UserIds: uids,
	})
	if err != nil {
		return nil, err
	}

	usm := make(map[string]*accountrpc.User)
	for _, v := range users.List {
		usm[v.UserId] = v
	}

	var list []*types.RemarkBackDTO
	for _, v := range out.List {
		m := ConvertRemarkTypes(v, usm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
