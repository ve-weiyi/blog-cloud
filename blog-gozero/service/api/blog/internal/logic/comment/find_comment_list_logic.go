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

func (l *FindCommentListLogic) FindCommentList(req *types.QueryCommentReq) (resp *types.PageResp, err error) {
	in := &messagerpc.FindCommentReplyListReq{
		Paginate: &messagerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TopicId:  req.TopicId,
		ParentId: req.ParentId,
		ReplyId:  0,
		Type:     req.Type,
	}

	// 查找评论列表
	out, err := l.svcCtx.MessageRpc.FindCommentReplyList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	usm, err := apiutils.BatchQueryMulti(out.List,
		func(v *messagerpc.Comment) []string {
			return []string{v.UserId, v.ReplyUserId}
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
		func(v *messagerpc.Comment) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	// 查找评论回复列表
	list := make([]*types.Comment, 0)
	for _, v := range out.List {
		m := &types.Comment{
			Id:               v.Id,
			UserId:           v.UserId,
			TerminalId:       v.TerminalId,
			TopicId:          v.TopicId,
			ParentId:         v.ParentId,
			ReplyId:          v.ReplyId,
			ReplyUserId:      v.ReplyUserId,
			CommentContent:   v.CommentContent,
			Type:             v.Type,
			CreatedAt:        v.CreatedAt,
			LikeCount:        v.LikeCount,
			ClientInfo:       vsm[v.TerminalId],
			UserInfo:         usm[v.UserId],
			ReplyUserInfo:    usm[v.ReplyUserId],
			ReplyCount:       0,
			CommentReplyList: make([]*types.CommentReply, 0),
		}

		// 查询回复评论
		reply, _ := l.svcCtx.MessageRpc.FindCommentReplyList(l.ctx, &messagerpc.FindCommentReplyListReq{
			Paginate: &messagerpc.PageReq{
				Page:     1,
				PageSize: 3,
				Sorts:    []string{"created_at desc"},
			},
			TopicId:  req.TopicId,
			ParentId: v.Id,
			ReplyId:  0,
			Type:     req.Type,
		})

		for _, r := range reply.List {
			mr := &types.CommentReply{
				Id:             r.Id,
				UserId:         r.UserId,
				TerminalId:     r.TerminalId,
				TopicId:        r.TopicId,
				ParentId:       r.ParentId,
				ReplyId:        r.ReplyId,
				ReplyUserId:    r.ReplyUserId,
				CommentContent: r.CommentContent,
				Status:         r.Status,
				Type:           r.Type,
				CreatedAt:      r.CreatedAt,
				LikeCount:      r.LikeCount,
				ClientInfo:     vsm[r.TerminalId],
				UserInfo:       usm[r.UserId],
				ReplyUserInfo:  usm[r.ReplyUserId],
			}
			m.CommentReplyList = append(m.CommentReplyList, mr)
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
