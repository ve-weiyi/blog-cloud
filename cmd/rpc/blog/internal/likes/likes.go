// Package likes 以 t_like 为唯一状态：是否已赞与点赞数都在读取时派生。
package likes

import (
	"context"

	"gorm.io/gorm"

	"github.com/ve-weiyi/blog-cloud/model"
)

// Deps 是点赞模块的依赖。
type Deps struct {
	DB    *gorm.DB
	Likes model.TLikeModel
}

// Like 点赞，返回本次调用是否改变了状态；重复点赞返回 false。
//
// 去重由唯一键承担：冲突时插入不修改任何列，受影响行数为 0。
func Like(ctx context.Context, d Deps, targetType, targetID int64, userID string) (bool, error) {
	rows, err := d.Likes.InsertIgnore(ctx, &model.TLike{
		TargetType: targetType,
		TargetId:   targetID,
		UserId:     userID,
	})
	return rows == 1, err
}

// Unlike 取消点赞，返回本次调用是否改变了状态；未点赞返回 false。
func Unlike(ctx context.Context, d Deps, targetType, targetID int64, userID string) (bool, error) {
	rows, err := d.Likes.DeleteBatch(ctx,
		"target_type = ? AND target_id = ? AND user_id = ?", targetType, targetID, userID)
	return rows == 1, err
}

// Liked 返回目标 ID 到「当前用户是否已赞」的映射；未登录时全部为未赞。
func Liked(ctx context.Context, d Deps, targetType int64, userID string, ids []int64) (map[int64]bool, error) {
	liked := make(map[int64]bool, len(ids))
	if userID == "" || len(ids) == 0 {
		return liked, nil
	}

	list, err := d.Likes.FindAll(ctx,
		"user_id = ? AND target_type = ? AND target_id IN ?", userID, targetType, ids)
	if err != nil {
		return nil, err
	}
	for _, row := range list {
		liked[row.TargetId] = true
	}
	return liked, nil
}

// Counts 返回目标 ID 到点赞数的映射；无人点赞的目标不出现在结果里。
// 分组计数属于聚合查询，不在通用 CRUD 契约内，因此直接按模型表名下发。
func Counts(ctx context.Context, d Deps, targetType int64, ids []int64) (map[int64]int64, error) {
	counts := make(map[int64]int64, len(ids))
	if len(ids) == 0 {
		return counts, nil
	}

	var rows []struct {
		TargetId int64
		Total    int64
	}
	if err := d.DB.WithContext(ctx).Table(d.Likes.TableName()).
		Select("target_id, COUNT(*) AS total").
		Where("target_type = ? AND target_id IN ?", targetType, ids).
		Group("target_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		counts[row.TargetId] = row.Total
	}
	return counts, nil
}
