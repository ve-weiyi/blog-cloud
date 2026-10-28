package discussionservicelogic

import (
	"context"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/likes"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/grpc/metax"
	"github.com/ve-weiyi/blog-cloud/model"
	"github.com/ve-weiyi/vkit/x/jsonv"
)

func convertMessageOut(in *model.TMessage) *discussionrpc.Message {
	return &discussionrpc.Message{
		Id:             in.Id,
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		MessageContent: in.MessageContent,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
	}
}

func convertCommentOut(in *model.TComment, liked map[int64]bool, counts map[int64]int64) *discussionrpc.Comment {
	return &discussionrpc.Comment{
		Id:             in.Id,
		UserId:         in.UserId,
		DeviceId:       in.DeviceId,
		TopicId:        in.TopicId,
		ParentId:       in.ParentId,
		ReplyId:        in.ReplyId,
		ReplyUserId:    in.ReplyUserId,
		CommentContent: in.CommentContent,
		Type:           in.Type,
		Status:         in.Status,
		IsReview:       in.IsReview,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
		LikeCount:      counts[in.Id],
		Liked:          liked[in.Id],
	}
}

func convertTalkOut(in *model.TTalk, liked map[int64]bool, counts map[int64]int64) *discussionrpc.Talk {
	var images []string
	if in.Images != nil {
		jsonv.JsonToAny(*in.Images, &images)
	}

	return &discussionrpc.Talk{
		Id:        in.Id,
		UserId:    in.UserId,
		Content:   in.Content,
		Images:    images,
		IsTop:     in.IsTop,
		Status:    in.Status,
		CreatedAt: in.CreatedAt.UnixMilli(),
		UpdatedAt: in.UpdatedAt.UnixMilli(),
		LikeCount: counts[in.Id],
		Liked:     liked[in.Id],
	}
}

// likeStats 批量取当前用户的点赞状态与各目标的点赞数。
func likeStats(ctx context.Context, svcCtx *svc.ServiceContext, targetType int64, ids []int64) (liked map[int64]bool, counts map[int64]int64, err error) {
	uid, _ := metax.GetUserIdFromCtx(ctx)
	deps := svcCtx.LikeDeps()

	liked, err = likes.Liked(ctx, deps, targetType, uid, ids)
	if err != nil {
		return nil, nil, err
	}
	counts, err = likes.Counts(ctx, deps, targetType, ids)
	if err != nil {
		return nil, nil, err
	}
	return liked, counts, nil
}

// commentIDs 提取评论 ID，供批量查询使用。
func commentIDs(records []*model.TComment) []int64 {
	ids := make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.Id)
	}
	return ids
}
