package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentReplyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCommentReplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentReplyListLogic {
	return &FindCommentReplyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询评论回复列表
func (l *FindCommentReplyListLogic) FindCommentReplyList(in *messagerpc.FindCommentReplyListReq) (*messagerpc.FindCommentReplyListResp, error) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.Type != 0 {
		opts = append(opts, query.WithCondition("type = ?", in.Type))
	}

	if in.TopicId != 0 {
		opts = append(opts, query.WithCondition("topic_id = ?", in.TopicId))
	}

	if in.ParentId >= 0 {
		opts = append(opts, query.WithCondition("parent_id = ?", in.ParentId))
	}
	page, size, sorts, conditions, params := query.NewQueryBuilder(opts...).Build()

	records, total, err := l.svcCtx.TCommentModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*messagerpc.CommentDetailsResp
	for _, v := range records {
		list = append(list, convertCommentOut(v))
	}

	return &messagerpc.FindCommentReplyListResp{
		List: list,
		Pagination: &messagerpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}
