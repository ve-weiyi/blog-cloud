package discussionservicelogic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/discussionrpc"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/svc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/enums"
	"github.com/ve-weiyi/blog-cloud/model"
)

type ListCommentRepliesByParentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentRepliesByParentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentRepliesByParentsLogic {
	return &ListCommentRepliesByParentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量查询多个父评论的回复，每个父评论最多返回 limit_per_parent 条。
func (l *ListCommentRepliesByParentsLogic) ListCommentRepliesByParents(in *discussionrpc.ListCommentRepliesByParentsRequest) (*discussionrpc.ListCommentRepliesByParentsResponse, error) {
	if len(in.ParentIds) == 0 {
		return &discussionrpc.ListCommentRepliesByParentsResponse{}, nil
	}

	limit := in.LimitPerParent
	if limit <= 0 {
		limit = 3
	}

	conditions := "`parent_id` IN ?"
	args := []any{in.ParentIds}
	if in.IsReview != nil {
		conditions += " AND `is_review` = ?"
		args = append(args, *in.IsReview)
	}
	if in.ExcludeDeleted {
		conditions += " AND `status` != ?"
		args = append(args, enums.CommentStatusDeleted)
	}
	args = append(args, limit)

	query := fmt.Sprintf("SELECT * FROM ("+
		"SELECT `t_comment`.*, ROW_NUMBER() OVER (PARTITION BY `parent_id` ORDER BY `id` ASC) AS rn "+
		"FROM `t_comment` WHERE %s) t WHERE rn <= ? ORDER BY `parent_id` ASC, `id` ASC", conditions)

	var records []*model.TComment
	if err := l.svcCtx.GormDB.Raw(query, args...).Scan(&records).Error; err != nil {
		return nil, err
	}

	var countRows []struct {
		ParentId int64 `gorm:"column:parent_id"`
		Total    int64 `gorm:"column:total"`
	}
	if err := l.svcCtx.GormDB.Table("`t_comment`").
		Select("parent_id, COUNT(*) AS total").
		Where(conditions, args[:len(args)-1]...).
		Group("parent_id").
		Scan(&countRows).Error; err != nil {
		return nil, err
	}
	totals := make(map[int64]int64, len(countRows))
	for _, row := range countRows {
		totals[row.ParentId] = row.Total
	}

	list := make([]*discussionrpc.Comment, 0, len(records))
	lm, lcm, err := likeStats(l.ctx, l.svcCtx, enums.LikeTargetComment, commentIDs(records))
	if err != nil {
		return nil, err
	}

	for _, entity := range records {
		list = append(list, convertCommentOut(entity, lm, lcm))
	}

	return &discussionrpc.ListCommentRepliesByParentsResponse{
		List:         list,
		ParentTotals: totals,
	}, nil
}
