package cachekey

import "fmt"

// 网站日访客集合
func GetDailyUserVisitKey(day string) string {
	return fmt.Sprintf("blog:visit:visitor:%v", day)
}

const (
	// 文章访问量排行
	ArticleViewCountKey = "blog:article:view_count"
	// 文章分块重建的增量水位（updated_at 秒 + id）
	ArticleChunkWatermarkKey = "blog:article:chunk:watermark"
	// 网站日 UV
	DailyUserViewCountKey = "blog:visit:daily_uv"
	// 网站日 PV
	DailyPageViewCountKey = "blog:visit:daily_pv"
	// 总 UV
	TotalUserViewCountKey = "blog:visit:total_uv"
	// 总 PV
	TotalPageViewCountKey = "blog:visit:total_pv"
)
