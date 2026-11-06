package cachekey

const (
	// 限流器前缀
	RateLimitStrictPrefix = "blog:ratelimit:strict:"
	// RateLimitMusicPrefix 音乐代理的独立额度：播放器一次加载会并发打多个接口
	RateLimitMusicPrefix = "blog:ratelimit:music:"
)
