package cache

import (
	"fmt"
	"time"
)

const (
	ExpireTimeMinute = time.Minute
	ExpireTimeHour   = time.Hour
	ExpireTimeDay    = 24 * time.Hour
	ExpireTimeWeek   = 7 * 24 * time.Hour
)

const (
	ArticleUserLike  = "article_user_like"
	ArticleLikeCount = "article_like_count"
)

const (
	CommentUserLike  = "comment_user_like"
	CommentLikeCount = "comment_like_count"
)

const (
	TalkUserLike  = "talk_user_like"
	TalkLikeCount = "talk_like_count"
)

const (
	UserOnline = "user_online"
)

func WrapCacheKey(prefix string, key interface{}) string {
	return fmt.Sprintf("%s:%v", prefix, key)
}
