package commentrpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/pb/commentrpc"
	"github.com/ve-weiyi/ve-blog-golang/zero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindTopicCommentCountsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindTopicCommentCountsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindTopicCommentCountsLogic {
	return &FindTopicCommentCountsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论回复数量
func (l *FindTopicCommentCountsLogic) FindTopicCommentCounts(in *commentrpc.IdsReq) (*commentrpc.FindTopicCommentCountsResp, error) {

	cm := make(map[int64]int64)
	for _, v := range in.Ids {
		count, err := l.svcCtx.TCommentModel.FindCount(l.ctx, "topic_id=?", v)
		if err != nil {
			return nil, err
		}

		cm[v] = count
	}

	return &commentrpc.FindTopicCommentCountsResp{
		TopicCommentCounts: cm,
	}, nil
}
