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

type QueryRecentCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取最新评论列表
func NewQueryRecentCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRecentCommentListLogic {
	return &QueryRecentCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRecentCommentListLogic) QueryRecentCommentList(req *types.QueryCommentListReq) (resp *types.ListResult, err error) {
	isReview := int64(enums.CommentIsReviewYes)

	out, err := l.svcCtx.DiscussionService.ListCommentReplies(l.ctx, &discussionservice.ListCommentRepliesRequest{
		ListQuery: &discussionservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TopicId:        req.TopicId,
		Type:           req.Type,
		IsReview:       &isReview,
		ExcludeDeleted: true,
	})
	if err != nil {
		return nil, err
	}

	usm, err := enrich.GetUserInfos(l.ctx, l.svcCtx,
		collectx.DistinctMany(out.List, func(v *discussionservice.Comment) []string {
			return []string{v.UserId, v.ReplyUserId}
		}),
	)
	if err != nil {
		return nil, err
	}

	vsm, err := enrich.GetGuests(l.ctx, l.svcCtx,
		collectx.Distinct(out.List, func(v *discussionservice.Comment) string {
			return v.DeviceId
		}),
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Comment, 0)
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
			GuestInfo:        vsm[v.DeviceId],
			UserInfo:         usm[v.UserId],
			ReplyUserInfo:    usm[v.ReplyUserId],
			ReplyCount:       0,
			CommentReplyList: make([]*types.CommentReply, 0),
		})
	}

	return &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}, nil
}
