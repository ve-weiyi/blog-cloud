package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/accountrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/client/commentrpc"
)

type FindCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询评论列表
func NewFindCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentListLogic {
	return &FindCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindCommentListLogic) FindCommentList(req *types.CommentQueryReq) (resp *types.PageResp, err error) {
	in := &commentrpc.FindCommentListReq{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Sorts:      req.Sorts,
		TopicId:    req.TopicId,
		ParentId:   req.ParentId,
		ReplyMsgId: 0,
		Type:       req.Type,
	}

	// 查找评论列表
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

	// 查找评论回复列表
	list := make([]*types.Comment, 0)
	for _, v := range out.List {
		m := ConvertCommentTypes(v, usm)
		// 查询回复评论
		reply, _ := l.svcCtx.CommentRpc.FindCommentReplyList(l.ctx, &commentrpc.FindCommentReplyListReq{
			Page:       1,
			PageSize:   3,
			Sorts:      in.Sorts,
			TopicId:    req.TopicId,
			ParentId:   v.Id,
			ReplyMsgId: 0,
			Type:       req.Type,
		})

		for _, r := range reply.List {
			m.CommentReplyList = append(m.CommentReplyList, ConvertCommentReplyTypes(r, usm))
		}
		m.ReplyCount = reply.Total
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = in.Page
	resp.PageSize = in.PageSize
	resp.Total = out.Total
	resp.List = list
	return resp, nil
}

func ConvertCommentTypes(in *commentrpc.CommentDetails, usm map[string]*accountrpc.User) (out *types.Comment) {
	out = &types.Comment{
		Id:               in.Id,
		TopicId:          in.TopicId,
		ParentId:         in.ParentId,
		ReplyMsgId:       in.ReplyMsgId,
		UserId:           in.UserId,
		ReplyUserId:      in.ReplyUserId,
		CommentContent:   in.CommentContent,
		Type:             in.Type,
		CreatedAt:        in.CreatedAt,
		LikeCount:        in.LikeCount,
		ReplyCount:       0,
		CommentReplyList: make([]*types.CommentReply, 0),
	}

	// 用户信息
	if out.UserId != "" {
		user, ok := usm[out.UserId]
		if ok && user != nil {
			out.User = ConvertCommentUserInfoToPb(user)
		}
	}
	// 回复用户信息
	if out.ReplyUserId != "" {
		user, ok := usm[out.ReplyUserId]
		if ok && user != nil {
			out.ReplyUser = ConvertCommentUserInfoToPb(user)
		}
	}

	return
}

func ConvertCommentReplyTypes(req *commentrpc.CommentDetails, usm map[string]*accountrpc.User) (out *types.CommentReply) {
	out = &types.CommentReply{
		Id:             req.Id,
		TopicId:        req.TopicId,
		ParentId:       req.ParentId,
		ReplyMsgId:     req.ReplyMsgId,
		UserId:         req.UserId,
		ReplyUserId:    req.ReplyUserId,
		CommentContent: req.CommentContent,
		Type:           req.Type,
		CreatedAt:      req.CreatedAt,
		LikeCount:      req.LikeCount,
	}

	// 用户信息
	if out.UserId != "" {
		user, ok := usm[out.UserId]
		if ok && user != nil {
			out.User = ConvertCommentUserInfoToPb(user)
		}
	}
	// 回复用户信息
	if out.ReplyUserId != "" {
		user, ok := usm[out.ReplyUserId]
		if ok && user != nil {
			out.ReplyUser = ConvertCommentUserInfoToPb(user)
		}
	}

	return
}

func ConvertCommentUserInfoToPb(in *accountrpc.User) (out *types.CommentUserInfo) {
	return &types.CommentUserInfo{
		UserId:   in.UserId,
		Nickname: in.Nickname,
		Avatar:   in.Avatar,
		Website:  in.Info,
	}
}
