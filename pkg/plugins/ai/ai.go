package ai

import (
	"net/http"
	"strings"

	"github.com/openai/openai-go/v3/option"

	"github.com/ve-weiyi/ve-blog-golang/pkg/plugins/ai/chat"
)

type AiPlugin struct {
	proxy *chat.OpenAIProxy
}

func NewAiPlugin(opts ...option.RequestOption) *AiPlugin {
	proxy := chat.NewOpenAIProxy(opts...)
	return &AiPlugin{
		proxy: proxy,
	}
}

func (s *AiPlugin) Handler(prefix string) http.HandlerFunc {
	proxy := s.proxy
	return func(w http.ResponseWriter, r *http.Request) {
		// 仅允许 POST 请求
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, prefix)
		path = strings.TrimPrefix(path, "/")
		switch path {
		case "chat/completions":
			proxy.ChatCompletionsHandler(w, r)
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	}
}
