package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/blog/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
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

func (l *FindCommentRecentListLogic) FindCommentRecentList(req *types.QueryCommentReq) (resp *types.PageResp, err error) {
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
	out, err := l.svcCtx.MessageRpc.FindCommentReplyList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 查询用户信息
	usm, err := apiutils.BatchQueryMulti(out.List,
		func(v *messagerpc.CommentDetailsResp) []string {
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
		func(v *messagerpc.CommentDetailsResp) string {
			return v.TerminalId
		},
		func(ids []string) (map[string]*types.ClientInfoVO, error) {
			return apiutils.GetVisitorInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

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
			Status:           v.Status,
			Type:             v.Type,
			CreatedAt:        v.CreatedAt,
			LikeCount:        v.LikeCount,
			ClientInfo:       vsm[v.TerminalId],
			UserInfo:         usm[v.UserId],
			ReplyUserInfo:    usm[v.ReplyUserId],
			ReplyCount:       0,
			CommentReplyList: make([]*types.CommentReply, 0),
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
