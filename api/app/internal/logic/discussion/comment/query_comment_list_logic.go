package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
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
	status := int64(1)

	in := &discussionservice.ListCommentsRequest{
		ListQuery: &discussionservice.ListQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		Status: &status,
	}

	out, err := l.svcCtx.DiscussionService.ListComments(l.ctx, in)
	if err != nil {
		return nil, err
	}

	usm, err := apiutils.BatchQueryMulti(out.List,
		func(v *discussionservice.Comment) []string {
			return []string{v.UserId, v.ReplyUserId}
		},
		func(ids []string) (map[string]*types.UserInfoVO, error) {
			return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	vsm, err := apiutils.BatchQuery(out.List,
		func(v *discussionservice.Comment) string {
			return v.DeviceId
		},
		func(ids []string) (map[string]*types.GuestInfoVO, error) {
			return apiutils.GetGuests(l.ctx, l.svcCtx, ids)
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]*types.Comment, 0)
	for _, v := range out.List {
		pageSize := int64(3)
		replyOut, err := l.svcCtx.DiscussionService.ListCommentReplies(l.ctx, &discussionservice.ListCommentRepliesRequest{
			ListQuery: &discussionservice.ListQuery{
				Page:     1,
				PageSize: pageSize,
			},
			TopicId:  &v.TopicId,
			ParentId: &v.Id,
		})
		if err != nil {
			return nil, err
		}

		rusm, err := apiutils.BatchQueryMulti(replyOut.List,
			func(r *discussionservice.Comment) []string {
				return []string{r.UserId, r.ReplyUserId}
			},
			func(ids []string) (map[string]*types.UserInfoVO, error) {
				return apiutils.GetUserInfos(l.ctx, l.svcCtx, ids)
			},
		)
		if err != nil {
			return nil, err
		}

		rvsm, err := apiutils.BatchQuery(replyOut.List,
			func(r *discussionservice.Comment) string {
				return r.DeviceId
			},
			func(ids []string) (map[string]*types.GuestInfoVO, error) {
				return apiutils.GetGuests(l.ctx, l.svcCtx, ids)
			},
		)
		if err != nil {
			return nil, err
		}

		replies := make([]*types.CommentReply, 0)
		for _, r := range replyOut.List {
			replies = append(replies, &types.CommentReply{
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
				GuestInfo:      rvsm[r.DeviceId],
				UserInfo:       rusm[r.UserId],
				ReplyUserInfo:  rusm[r.ReplyUserId],
			})
		}

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
			ReplyCount:       replyOut.ListResult.Total,
			CommentReplyList: replies,
		})
	}

	resp = &types.ListResult{
		Page:     out.ListResult.Page,
		PageSize: out.ListResult.PageSize,
		Total:    out.ListResult.Total,
		List:     list,
	}
	return
}
