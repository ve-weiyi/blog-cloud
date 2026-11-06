package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/ve-weiyi/vkit/adapter/storagex"
)

type Config struct {
	rest.RestConf
	RedisConf     RedisConf
	StorageConfig storagex.StorageConfig
	TokenConf     TokenConf
	AiConf        AiConf
	// SiteURL 读者端站点根地址，用于拼接 AI 引用链接；为空时退化为相对路径
	SiteURL string `json:"site_url,optional"`

	AppRpcConf zrpc.RpcClientConf
}

// AiConf 是 AI 助手的模型接入、上下文预算与配额配置。
type AiConf struct {
	BaseURL         string          `json:"base_url,optional"`                       // OpenAI 兼容端点，如 https://api.deepseek.com/v1
	APIKey          string          `json:"api_key,optional"`                        // 留空时读环境变量 AI_API_KEY
	Model           string          `json:"model,optional,default=deepseek-chat"`    // 模型名
	Temperature     float64         `json:"temperature,optional,default=0.3"`        // 采样温度
	MaxTokens       int             `json:"max_tokens,optional,default=800"`         // 单次回答的输出上限
	Timeout         time.Duration   `json:"timeout,optional,default=60s"`            // 上游调用超时
	MaxHistoryTurns int             `json:"max_history_turns,optional,default=6"`    // 保留的历史轮次
	MaxHistoryChars int             `json:"max_history_chars,optional,default=4000"` // 历史正文总字符上限
	MaxPromptChars  int             `json:"max_prompt_chars,optional,default=12000"` // 提示词总字符上限
	SwitchTTL       time.Duration   `json:"switch_ttl,optional,default=60s"`         // 站点开关缓存时长
	Retrieval       AiRetrievalConf `json:"retrieval,optional"`                      // 检索参数
	Quota           AiQuotaConf     `json:"quota,optional"`                          // 配额参数
}

// AiRetrievalConf 是检索参数。
type AiRetrievalConf struct {
	TopK          int `json:"top_k,optional,default=5"`              // 片段数量
	MaxCharsTotal int `json:"max_chars_total,optional,default=4000"` // 资料总字符上限
}

// AiQuotaConf 是配额参数。
type AiQuotaConf struct {
	GuestDaily    int64         `json:"guest_daily,optional,default=20"`    // 游客每日次数
	UserDaily     int64         `json:"user_daily,optional,default=100"`    // 登录用户每日次数
	IPDaily       int64         `json:"ip_daily,optional,default=200"`      // 单来源每日次数
	MaxConcurrent int           `json:"max_concurrent,optional,default=32"` // 全站并发流上限
	Lease         time.Duration `json:"lease,optional,default=5m"`          // 并发租约，应大于 Timeout；进程崩溃后过期自愈
}

type TokenConf struct {
	AccessTokenTTL     time.Duration `json:"access_token_ttl,optional,default=2h"`
	RefreshTokenTTL    time.Duration `json:"refresh_token_ttl,optional,default=168h"` // idle
	RefreshAbsoluteTTL time.Duration `json:"refresh_absolute_ttl,optional,default=720h"`
}

// redis缓存配置
type RedisConf struct {
	Host     string `json:"host,default=127.0.0.1"`
	Port     string `json:"port,default=6379"`
	Password string `json:"password,default=''"`
	DB       int    `json:"db,default=0"`
}
