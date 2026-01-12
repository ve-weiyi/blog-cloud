package messagerpclogic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/model"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/common/query"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/pb/messagerpc"
	"github.com/ve-weiyi/ve-blog-golang/blog-gozero/service/rpc/blog/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindCommentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindCommentListLogic {
	return &FindCommentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页获取评论列表
func (l *FindCommentListLogic) FindCommentList(in *messagerpc.FindCommentListReq) (*messagerpc.FindCommentListResp, error) {
	var opts []query.Option
	if in.Paginate != nil {
		opts = append(opts, query.WithPage(int(in.Paginate.Page)))
		opts = append(opts, query.WithSize(int(in.Paginate.PageSize)))
		opts = append(opts, query.WithSorts(in.Paginate.Sorts...))
	}

	if in.UserId != "" {
		opts = append(opts, query.WithCondition("user_id = ?", in.UserId))
	}

	if in.Status >= 0 {
		opts = append(opts, query.WithCondition("status = ?", in.Status))
	}

	if in.Type != 0 {
		opts = append(opts, query.WithCondition("type = ?", in.Type))
	}

	page, size, sorts, conditions, params := query.NewQueryBuilder(opts...).Build()

	records, total, err := l.svcCtx.TCommentModel.FindListAndTotal(l.ctx, page, size, sorts, conditions, params...)
	if err != nil {
		return nil, err
	}

	var list []*messagerpc.Comment
	for _, v := range records {
		m := convertCommentOut(v)
		list = append(list, m)
	}

	return &messagerpc.FindCommentListResp{
		List: list,
		Pagination: &messagerpc.PageResp{
			Page:     int64(page),
			PageSize: int64(size),
			Total:    total,
		},
	}, nil
}

func convertCommentOut(in *model.TComment) (out *messagerpc.Comment) {
	out = &messagerpc.Comment{
		Id:             in.Id,
		UserId:         in.UserId,
		TopicId:        in.TopicId,
		ParentId:       in.ParentId,
		ReplyId:        in.ReplyId,
		ReplyUserId:    in.ReplyUserId,
		CommentContent: in.CommentContent,
		Type:           in.Type,
		Status:         in.Status,
		CreatedAt:      in.CreatedAt.UnixMilli(),
		UpdatedAt:      in.UpdatedAt.UnixMilli(),
		LikeCount:      in.LikeCount,
	}

	return out
}
