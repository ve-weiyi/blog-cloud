package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
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
	in := &messagerpc.FindCommentListReq{
		Paginate: &messagerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TopicId:    req.TopicId,
		ParentId:   req.ParentId,
		ReplyMsgId: 0,
		Type:       req.Type,
	}

	// 查找评论列表
	out, err := l.svcCtx.MessageRpc.FindCommentList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	var uids []string
	for _, v := range out.List {
		uids = append(uids, v.UserId)
		uids = append(uids, v.ReplyUserId)
	}

	// 查询用户信息
	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	// 查找评论回复列表
	list := make([]*types.Comment, 0)
	for _, v := range out.List {
		m := ConvertCommentTypes(v, usm)
		// 查询回复评论
		reply, _ := l.svcCtx.MessageRpc.FindCommentReplyList(l.ctx, &messagerpc.FindCommentReplyListReq{
			Paginate: &messagerpc.PageReq{
				Page:     1,
				PageSize: 3,
				Sorts:    []string{"created_at desc"},
			},
			TopicId:    req.TopicId,
			ParentId:   v.Id,
			ReplyMsgId: 0,
			Type:       req.Type,
		})

		for _, r := range reply.List {
			m.CommentReplyList = append(m.CommentReplyList, ConvertCommentReplyTypes(r, usm))
		}
		m.ReplyCount = reply.Pagination.Total
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}

func ConvertCommentTypes(in *messagerpc.CommentDetailsResp, usm map[string]*types.UserInfoVO) (out *types.Comment) {
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
			out.User = user
		}
	}
	// 回复用户信息
	if out.ReplyUserId != "" {
		user, ok := usm[out.ReplyUserId]
		if ok && user != nil {
			out.ReplyUser = user
		}
	}

	return
}

func ConvertCommentReplyTypes(req *messagerpc.CommentDetailsResp, usm map[string]*types.UserInfoVO) (out *types.CommentReply) {
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
			out.User = user
		}
	}
	// 回复用户信息
	if out.ReplyUserId != "" {
		user, ok := usm[out.ReplyUserId]
		if ok && user != nil {
			out.ReplyUser = user
		}
	}

	return
}
