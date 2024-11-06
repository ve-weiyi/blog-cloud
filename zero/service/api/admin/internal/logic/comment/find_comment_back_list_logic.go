package comment

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentBackListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询评论列表(后台)
func NewFindCommentBackListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentBackListLogic {
	return &FindCommentBackListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindCommentBackListLogic) FindCommentBackList(req *types.CommentQuery) (resp *types.PageResp, err error) {
	in := &commentrpc.FindCommentListReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Sorts:      req.Sorts,
		ReplyMsgId: 0,
		Type:       req.Type,
	}

	// 查找评论列表
	out, err := l.svcCtx.CommentRpc.FindCommentList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	var aids []int64
	for _, v := range out.List {
		uids = append(uids, v.UserId)
		uids = append(uids, v.ReplyUserId)
		aids = append(aids, v.TopicId)
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

	// 查询文章信息
	topics, err := l.svcCtx.ArticleRpc.FindArticlePreviewList(l.ctx, &articlerpc.FindArticlePreviewListReq{
		Ids: aids,
	})

	tsm := make(map[int64]*articlerpc.ArticlePreview)
	for _, v := range topics.List {
		tsm[v.Id] = v
	}

	// 查找评论回复列表
	var list []*types.CommentBackDTO
	for _, v := range out.List {
		m := ConvertCommentTypes(v, usm, tsm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertCommentTypes(in *commentrpc.CommentDetails, usm map[string]*accountrpc.User, tsm map[int64]*articlerpc.ArticlePreview) (out *types.CommentBackDTO) {
	out = &types.CommentBackDTO{
		Id:             in.Id,
		Type:           in.Type,
		TopicTitle:     "",
		Avatar:         "",
		Nickname:       "",
		ToNickname:     "",
		CommentContent: in.CommentContent,
		IsReview:       in.IsReview,
		CreatedAt:      in.CreatedAt,
	}

	// 文章信息
	if in.TopicId != 0 {
		topic, ok := tsm[in.TopicId]
		if ok && topic != nil {
			out.TopicTitle = topic.ArticleTitle
		}
	}

	// 用户信息
	if in.UserId != "" {
		user, ok := usm[in.UserId]
		if ok && user != nil {
			out.Nickname = user.Nickname
			out.Avatar = user.Avatar
		}
	}

	if in.ReplyUserId != "" {
		user, ok := usm[in.ReplyUserId]
		if ok && user != nil {
			out.ToNickname = user.Nickname
		}
	}

	return
}
