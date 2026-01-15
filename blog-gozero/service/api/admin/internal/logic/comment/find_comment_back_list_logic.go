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

func (l *FindCommentBackListLogic) FindCommentBackList(req *types.QueryCommentReq) (resp *types.PageResp, err error) {
	in := &messagerpc.FindCommentListReq{
		Paginate: &messagerpc.PageReq{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		UserId: req.UserId,
		Status: req.Status,
		Type:   req.Type,
	}

	// 查找评论列表
	out, err := l.svcCtx.MessageRpc.FindCommentList(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 提取用户ID和文章ID
	var aids []int64
	for _, v := range out.List {
		aids = append(aids, v.TopicId)
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

	// 查询文章信息
	topics, err := l.svcCtx.ArticleRpc.FindArticlePreviewList(l.ctx, &articlerpc.FindArticleListReq{
		Ids: aids,
	})
	if err != nil {
		return nil, err
	}

	tsm := make(map[int64]*articlerpc.ArticlePreview)
	for _, v := range topics.List {
		tsm[v.Id] = v
	}

	// 查找评论回复列表
	var list []*types.CommentBackVO
	for _, v := range out.List {
		list = append(list, &types.CommentBackVO{
			Id:         v.Id,
			UserId:     v.UserId,
			TerminalId: v.TerminalId,
			Type:       v.Type,
			TopicTitle: func() string {
				if t := tsm[v.TopicId]; t != nil {
					return t.ArticleTitle
				}
				return ""
			}(),
			ReplyUserId:    v.ReplyUserId,
			CommentContent: v.CommentContent,
			Status:         v.Status,
			CreatedAt:      v.CreatedAt,
			UserInfo:       usm[v.UserId],
			ClientInfo:     vsm[v.TerminalId],
			ReplyUserInfo:  usm[v.ReplyUserId],
		})
	}

	resp = &types.PageResp{}
	resp.Page = out.Pagination.Page
	resp.PageSize = out.Pagination.PageSize
	resp.Total = out.Pagination.Total
	resp.List = list
	return resp, nil
}
