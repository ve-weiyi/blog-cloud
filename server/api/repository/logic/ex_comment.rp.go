package logic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/cache"
)

// 点赞评论
func (s *CommentRepository) LikeComment(ctx context.Context, uid int, commentId int) (data interface{}, err error) {
	// 用户点赞的评论列表
	commentUserLikeKey := cache.WrapCacheKey(cache.CommentUserLike, uid)
	// 当前评论的点赞量
	commentLikeCountKey := cache.WrapCacheKey(cache.CommentLikeCount, commentId)

	// 判断是否已经点赞
	if s.Cache.SIsMember(ctx, commentUserLikeKey, commentId).Val() {
		// 点过赞则删除评论id
		s.Cache.SRem(ctx, commentUserLikeKey, commentId)
		// 评论点赞量-1
		s.Cache.Decr(ctx, commentLikeCountKey)
	} else {
		// 未点赞则增加评论id
		s.Cache.SAdd(ctx, commentUserLikeKey, commentId)
		// 评论点赞量+1
		s.Cache.Incr(ctx, commentLikeCountKey)
	}

	return data, nil
}

// 获取用户点赞记录
func (s *CommentRepository) FindUserLikeComment(ctx context.Context, uid int) (data []string, err error) {
	// 用户点赞的评论列表
	commentUserLikeKey := cache.WrapCacheKey(cache.CommentUserLike, uid)
	return s.Cache.SMembers(ctx, commentUserLikeKey).Result()
}
