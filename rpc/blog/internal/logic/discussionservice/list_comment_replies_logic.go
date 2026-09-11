package discussionservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/rpc/blog/internal/svc"
	"github.com/ve-weiyi/vkit/adapter/gormx/queryx"
)

type ListCommentRepliesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentRepliesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentRepliesLogic {
	return &ListCommentRepliesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCommentRepliesLogic) ListCommentReplies(in *discussionrpc.ListCommentRepliesRequest) (*discussionrpc.ListCommentRepliesResponse, error) {
	var opts []queryx.Option
	if in.PageQuery != nil {
		opts = append(opts, queryx.WithPage(int(in.PageQuery.Page)))
		opts = append(opts, queryx.WithSize(int(in.PageQuery.PageSize)))
		opts = append(opts, queryx.WithSorts(in.PageQuery.Sorts...))
	}

	if in.TopicId != nil {
		opts = append(opts, queryx.WithCondition("topic_id = ?", *in.TopicId))
	}
	if in.ParentId != nil {
		opts = append(opts, queryx.WithCondition("parent_id = ?", *in.ParentId))
	}
	if in.ReplyId != nil {
		opts = append(opts, queryx.WithCondition("reply_id = ?", *in.ReplyId))
	}
	if in.Type != nil {
		opts = append(opts, queryx.WithCondition("type = ?", *in.Type))
	}

	page, size, sorts, conditions, params := queryx.NewQueryBuilder(opts...).Build()
	records, total, err := l.svcCtx.TCommentModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*discussionrpc.Comment
	for _, entity := range records {
		list = append(list, convertCommentOut(entity))
	}

	return &discussionrpc.ListCommentRepliesResponse{
		PageResult: &discussionrpc.PageResult{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
		List: list,
	}, nil
}
