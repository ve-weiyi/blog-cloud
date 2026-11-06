// ai-check 手动验证 AI 适配器能否连上真实的 OpenAI 兼容端点。
//
// 需要真实模型密钥，属非正式诊断程序：缺少 AI_API_KEY 时显式失败，不会回落到默认值。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ve-weiyi/blog-cloud/infra/ai"
)

func main() {
	baseURL := flag.String("base-url", envOr("AI_BASE_URL", "https://api.deepseek.com/v1"), "OpenAI 兼容端点")
	model := flag.String("model", envOr("AI_MODEL", "deepseek-chat"), "模型名")
	prompt := flag.String("prompt", "用一句话介绍你自己", "测试问题")
	timeout := flag.Duration("timeout", 60*time.Second, "整条流的总时限")
	flag.Parse()

	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "缺少环境变量 AI_API_KEY，拒绝运行")
		os.Exit(1)
	}

	client := ai.NewClient(ai.Config{
		BaseURL:   *baseURL,
		APIKey:    apiKey,
		Model:     *model,
		MaxTokens: 200,
		Timeout:   *timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	fmt.Printf("endpoint=%s model=%s prompt=%q\n", *baseURL, *model, *prompt)
	deltas := 0
	usage, err := client.Stream(ctx, []ai.Message{{Role: "user", Content: *prompt}}, func(delta string) error {
		deltas++
		fmt.Print(delta)
		return nil
	})
	fmt.Println()
	if err != nil {
		fmt.Fprintf(os.Stderr, "stream failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("usage: prompt=%d completion=%d deltas=%d\n", usage.PromptTokens, usage.CompletionTokens, deltas)
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
