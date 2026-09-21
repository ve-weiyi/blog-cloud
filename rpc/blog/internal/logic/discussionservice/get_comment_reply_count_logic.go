package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
)

type GetCommentReplyCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentReplyCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentReplyCountLogic {
	return &GetCommentReplyCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论回复数量
func (l *GetCommentReplyCountLogic) GetCommentReplyCount(in *discussionrpc.GetCommentReplyCountRequest) (*discussionrpc.GetCommentReplyCountResponse, error) {
	result := make(map[int64]int64, len(in.TopicIds))
	for _, topicId := range in.TopicIds {
		count, err := l.svcCtx.TCommentModel.FindCount(l.ctx, "topic_id = ?", topicId)
		if err != nil {
			return nil, err
		}
		result[topicId] = count
	}

	return &discussionrpc.GetCommentReplyCountResponse{TopicCommentCounts: result}, nil
}
