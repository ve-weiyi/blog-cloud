package friend

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/socialrpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/syslogrpc"
)

type FindFriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取友链列表
func NewFindFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindFriendListLogic {
	return &FindFriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindFriendListLogic) FindFriendList(req *types.QueryFriendReq) (resp *types.PageResp, err error) {
	in := &socialrpc.FindFriendListReq{
		Paginate: &socialrpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
	}
	out, err := l.svcCtx.SocialRpc.FindFriendList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Friend, 0)
	for _, v := range out.List {
		list = append(list, &types.Friend{
			Id:          v.Id,
			LinkName:    v.LinkName,
			LinkAvatar:  v.LinkAvatar,
			LinkAddress: v.LinkAddress,
			LinkIntro:   v.LinkIntro,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	_, err = l.svcCtx.SyslogRpc.AddVisitLog(l.ctx, &syslogrpc.AddVisitLogReq{
		PageName: "友链",
	})
	if err != nil {
		return nil, err
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
