package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/common/apiutils"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/svc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/api/admin/internal/types"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/articlerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/client/messagerpc"
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
	in := &messagerpc.FindCommentListReq{
		Paginate: &messagerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		ReplyMsgId: 0,
		Type:       req.Type,
	}

	// 查找评论列表
	out, err := l.svcCtx.MessageRpc.FindCommentList(l.ctx, in)
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
	usm, err := apiutils.GetUserInfos(l.ctx, l.svcCtx, uids)
	if err != nil {
		return nil, err
	}

	// 查询文章信息
	topics, err := l.svcCtx.ArticleRpc.FindArticlePreviewList(l.ctx, &articlerpc.FindArticleListReq{
		Ids: aids,
	})

	tsm := make(map[int64]*articlerpc.ArticlePreview)
	for _, v := range topics.List {
		tsm[v.Id] = v
	}

	// 查找评论回复列表
	var list []*types.CommentBackVO
	for _, v := range out.List {
		m := ConvertCommentTypes(v, usm, tsm)
		list = append(list, m)
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}

func ConvertCommentTypes(in *messagerpc.CommentDetailsResp, usm map[string]*types.UserInfoVO, tsm map[int64]*articlerpc.ArticlePreview) (out *types.CommentBackVO) {
	out = &types.CommentBackVO{
		Id:             in.Id,
		Type:           in.Type,
		TopicTitle:     "",
		UserId:         in.UserId,
		ReplyUserId:    in.ReplyUserId,
		CommentContent: in.CommentContent,
		IsReview:       in.IsReview,
		CreatedAt:      in.CreatedAt,
		User:           nil,
		ReplyUser:      nil,
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
			out.User = user
		}
	}

	if in.ReplyUserId != "" {
		user, ok := usm[in.ReplyUserId]
		if ok && user != nil {
			out.ReplyUser = user
		}
	}

	return
}
