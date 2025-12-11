package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentReplyCountsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCommentReplyCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentReplyCountsLogic {
	return &FindCommentReplyCountsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论回复数量
func (l *FindCommentReplyCountsLogic) FindCommentReplyCounts(in *messagerpc.FindCommentReplyCountsReq) (*messagerpc.FindCommentReplyCountsResp, error) {
	cm := make(map[int64]int64)
	for _, v := range in.TopicIds {
		count, err := l.svcCtx.TCommentModel.FindCount(l.ctx, "topic_id=?", v)
		if err != nil {
			return nil, err
		}

		cm[v] = count
	}

	return &messagerpc.FindCommentReplyCountsResp{
		TopicCommentCounts: cm,
	}, nil
}
