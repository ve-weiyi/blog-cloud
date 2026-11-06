package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/siteservice"
	"github.com/ve-weiyi/blog-cloud/infra/domain/configkey"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizcode"
	"github.com/ve-weiyi/blog-cloud/infra/errors/bizerr"
	"github.com/ve-weiyi/blog-cloud/infra/http/responsex"
	"github.com/ve-weiyi/vkit/x/jsonv"
)

// aiSwitchErrorTTL 读取开关失败后的短暂冷却，避免每请求都打 RPC。
const aiSwitchErrorTTL = 5 * time.Second

// AiSwitchMiddleware 按站点配置拦截 AI 请求：前端隐藏入口挡不住直接调用。
type AiSwitchMiddleware struct {
	siteService siteservice.SiteService
	ttl         time.Duration

	mu        sync.Mutex
	enabled   bool
	expiresAt time.Time
}

func NewAiSwitchMiddleware(siteService siteservice.SiteService, ttl time.Duration) *AiSwitchMiddleware {
	if ttl <= 0 {
		ttl = time.Minute
	}
	return &AiSwitchMiddleware{siteService: siteService, ttl: ttl}
}

func (m *AiSwitchMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.isEnabled(r.Context()) {
			responsex.Response(r, w, nil, bizerr.NewBizError(bizcode.CodeOperationNotAllowed, "AI 助手未开启"))
			return
		}
		next(w, r)
	}
}

func (m *AiSwitchMiddleware) isEnabled(ctx context.Context) bool {
	m.mu.Lock()
	if time.Now().Before(m.expiresAt) {
		enabled := m.enabled
		m.mu.Unlock()
		return enabled
	}
	m.mu.Unlock()

	enabled, ttl := m.load(ctx)

	m.mu.Lock()
	m.enabled = enabled
	m.expiresAt = time.Now().Add(ttl)
	m.mu.Unlock()
	return enabled
}

func (m *AiSwitchMiddleware) load(ctx context.Context) (bool, time.Duration) {
	out, err := m.siteService.GetConfig(ctx, &siteservice.GetConfigRequest{ConfigKey: configkey.Website})
	if err != nil {
		logx.WithContext(ctx).Errorf("读取 AI 开关失败: %v", err)
		return false, aiSwitchErrorTTL
	}

	var conf struct {
		WebsiteFeature struct {
			IsAiAssistant int64 `json:"is_ai_assistant"`
		} `json:"website_feature"`
	}
	if err := jsonv.JsonToAny(out.ConfigValue, &conf); err != nil {
		logx.WithContext(ctx).Errorf("解析 AI 开关失败: %v", err)
		return false, aiSwitchErrorTTL
	}
	return conf.WebsiteFeature.IsAiAssistant == 1, m.ttl
}
