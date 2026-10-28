package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/enrich"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/discussionservice"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/vkit/x/collectx"
)

type QueryCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取评论列表
func NewQueryCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCommentListLogic {
	return &QueryCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCommentListLogic) QueryCommentList(req *types.QueryCommentListReq) (resp *types.ListResult, err error) {
	isReview := int64(enums.CommentIsReviewYes)

	in := &discussionservice.ListCommentsRequest{
		ListQuery: &discussionservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TopicId:        req.TopicId,
		ParentId:       req.ParentId,
		Type:           req.Type,
		IsReview:       &isReview,
		ExcludeDeleted: true,
	}

	out, err := l.svcCtx.DiscussionService.ListComments(l.ctx, in)
	if err != nil {
		return nil, err
	}

	// 一次取回本页所有评论的前几条回复与回复总数
	parentIds := make([]int64, 0, len(out.List))
	for _, v := range out.List {
		parentIds = append(parentIds, v.Id)
	}

	var replyList []*discussionservice.Comment
	replyTotals := make(map[int64]int64)
	if len(parentIds) > 0 {
		replyOut, err := l.svcCtx.DiscussionService.ListCommentRepliesByParents(l.ctx, &discussionservice.ListCommentRepliesByParentsRequest{
			ParentIds:      parentIds,
			LimitPerParent: 3,
			IsReview:       &isReview,
			ExcludeDeleted: true,
		})
		if err != nil {
			return nil, err
		}
		replyList = replyOut.List
		replyTotals = replyOut.ParentTotals
	}

	// 评论与回复的用户/游客信息各批量查询一次
	allComments := make([]*discussionservice.Comment, 0, len(out.List)+len(replyList))
	allComments = append(allComments, out.List...)
	allComments = append(allComments, replyList...)

	usm, err := enrich.GetUserInfos(l.ctx, l.svcCtx,
		collectx.DistinctMany(allComments, func(v *discussionservice.Comment) []string {
			return []string{v.UserId, v.ReplyUserId}
		}),
	)
	if err != nil {
		return nil, err
	}

	vsm, err := enrich.GetGuests(l.ctx, l.svcCtx,
		collectx.Distinct(allComments, func(v *discussionservice.Comment) string {
			return v.DeviceId
		}),
	)
	if err != nil {
		return nil, err
	}

	repliesByParent := make(map[int64][]*types.CommentReply, len(parentIds))
	for _, r := range replyList {
		repliesByParent[r.ParentId] = append(repliesByParent[r.ParentId], &types.CommentReply{
			Id:             r.Id,
			UserId:         r.UserId,
			DeviceId:       r.DeviceId,
			TopicId:        r.TopicId,
			ParentId:       r.ParentId,
			ReplyId:        r.ReplyId,
			ReplyUserId:    r.ReplyUserId,
			CommentContent: r.CommentContent,
			Status:         r.Status,
			Type:           r.Type,
			CreatedAt:      r.CreatedAt,
			LikeCount:      r.LikeCount,
			Liked:          r.Liked,
			GuestInfo:      vsm[r.DeviceId],
			UserInfo:       usm[r.UserId],
			ReplyUserInfo:  usm[r.ReplyUserId],
		})
	}

	list := make([]*types.Comment, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, &types.Comment{
			Id:               v.Id,
			UserId:           v.UserId,
			DeviceId:         v.DeviceId,
			TopicId:          v.TopicId,
			ParentId:         v.ParentId,
			ReplyId:          v.ReplyId,
			ReplyUserId:      v.ReplyUserId,
			CommentContent:   v.CommentContent,
			Status:           v.Status,
			Type:             v.Type,
			CreatedAt:        v.CreatedAt,
			LikeCount:        v.LikeCount,
			Liked:            v.Liked,
			GuestInfo:        vsm[v.DeviceId],
			UserInfo:         usm[v.UserId],
			ReplyUserInfo:    usm[v.ReplyUserId],
			ReplyCount:       replyTotals[v.Id],
			CommentReplyList: repliesByParent[v.Id],
		})
	}

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		List:     list,
		Total:    out.ListResult.Total,
	}
	return
}
