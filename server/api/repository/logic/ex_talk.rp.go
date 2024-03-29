package logic

import (
	"context"

	"github.com/ve-weiyi/ve-blog-golang/server/infra/cache"
)

// 点赞说说
func (s *TalkRepository) LikeTalk(ctx context.Context, uid int, talkId int) (data interface{}, err error) {
	// 用户点赞的评论列表
	talkUserLikeKey := cache.WrapCacheKey(cache.TalkUserLike, uid)
	// 当前评论的点赞量
	talkLikeCountKey := cache.WrapCacheKey(cache.TalkLikeCount, talkId)

	// 判断是否已经点赞
	if s.Cache.SIsMember(ctx, talkUserLikeKey, talkId).Val() {
		// 点过赞则删除评论id
		s.Cache.SRem(ctx, talkUserLikeKey, talkId)
		// 评论点赞量-1
		s.Cache.Decr(ctx, talkLikeCountKey)
	} else {
		// 未点赞则增加评论id
		s.Cache.SAdd(ctx, talkUserLikeKey, talkId)
		// 评论点赞量+1
		s.Cache.Incr(ctx, talkLikeCountKey)
	}

	return data, nil
}

// 获取用户点赞记录
func (s *TalkRepository) FindUserLikeTalk(ctx context.Context, uid int) (data []string, err error) {
	// 用户点赞的评论列表
	talkUserLikeKey := cache.WrapCacheKey(cache.TalkUserLike, uid)
	return s.Cache.SMembers(ctx, talkUserLikeKey).Result()
}
