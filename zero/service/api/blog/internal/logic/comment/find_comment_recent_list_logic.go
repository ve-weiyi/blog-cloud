package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"
)

type FindCommentRecentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询最新评论回复列表
func NewFindCommentRecentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentRecentListLogic {
	return &FindCommentRecentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindCommentRecentListLogic) FindCommentRecentList(req *types.CommentQueryReq) (resp *types.PageResp, err error) {
	in := &commentrpc.FindCommentListReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Sorts:      req.Sorts,
		TopicId:    req.TopicId,
		ParentId:   req.ParentId,
		ReplyMsgId: 0,
		Type:       req.Type,
	}
	out, err := l.svcCtx.CommentRpc.FindCommentList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
		uids = append(uids, v.ReplyUserId)
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

	list := make([]*types.Comment, 0)
	for _, v := range out.List {
		m := ConvertCommentTypes(v, usm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}
