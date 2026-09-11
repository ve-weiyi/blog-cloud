package comment

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/api/app/internal/common/apiutils"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/client/discussionservice"
)

type QueryCommentReplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取评论回复列表
func NewQueryCommentReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryCommentReplyListLogic {
	return &QueryCommentReplyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryCommentReplyListLogic) QueryCommentReplyList(req *types.QueryCommentReplyListReq) (resp *types.PageResult, err error) {
	out, err := l.svcCtx.DiscussionService.ListCommentReplies(l.ctx, &discussionservice.ListCommentRepliesRequest{
		PageQuery: &discussionservice.PageQuery{
			Page:     req.Page,
			PageSize: req.PageSize,
			Sorts:    req.Sorts,
		},
		TopicId:  req.TopicId,
		ParentId: &req.CommentId,
		Type:     req.Type,
	})
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

	list := make([]*types.CommentReply, 0)
	for _, v := range out.List {
		list = append(list, &types.CommentReply{
			Id:             v.Id,
			UserId:         v.UserId,
			DeviceId:       v.DeviceId,
			TopicId:        v.TopicId,
			ParentId:       v.ParentId,
			ReplyId:        v.ReplyId,
			ReplyUserId:    v.ReplyUserId,
			CommentContent: v.CommentContent,
			Status:         v.Status,
			Type:           v.Type,
			CreatedAt:      v.CreatedAt,
			LikeCount:      v.LikeCount,
			GuestInfo:      vsm[v.DeviceId],
			UserInfo:       usm[v.UserId],
			ReplyUserInfo:  usm[v.ReplyUserId],
		})
	}

	return &types.PageResult{
		Page:     out.PageResult.Page,
		PageSize: out.PageResult.PageSize,
		Total:    out.PageResult.Total,
		List:     list,
	}, nil
}
