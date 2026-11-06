// Code scaffolded by goctl. Safe to edit.
package chat_stream

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aictx"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/ai/aiprompt"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/svc"
	"github.com/ve-weiyi/blog-cloud/cmd/api/app/internal/types"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/aiservice"
	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/client/contentservice"
	"github.com/ve-weiyi/blog-cloud/infra/ai"
)

const (
	// maxQuestionRunes 单轮问题的最大字符数。
	maxQuestionRunes = 1000
	// maxTitleRunes 会话标题取首个问题的前若干字符。
	maxTitleRunes = 20
	// noResultText 没有检索结果时的固定回答，不调用模型。
	noResultText = "站内没有找到相关内容。"
	// maxGuestHistoryMessages 游客一次最多携带的历史条数，超出从最旧的截断。
	maxGuestHistoryMessages = 50
	// maxGuestMessageRunes 游客每条历史的最大字符数。
	maxGuestMessageRunes = 2000
)

const (
	roleUser      = "user"
	roleAssistant = "assistant"
)

type AiChatStreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// AI 助手问答流（SSE）
func NewAiChatStreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AiChatStreamLogic {
	return &AiChatStreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AiChatStreamLogic) AiChatStream(req *types.AiChatStreamReq, client chan<- *types.AiStreamEvent) error {
	identity, ok := aictx.IdentityFrom(l.ctx)
	if !ok {
		return emit(l.ctx, client, aiprompt.Error("unauthenticated", "缺少身份信息"))
	}

	question := strings.TrimSpace(req.Message)
	if n := utf8.RuneCountInString(question); n == 0 || n > maxQuestionRunes {
		return emit(l.ctx, client, aiprompt.Error("invalid_param", fmt.Sprintf("消息长度必须在 1-%d 个字符之间", maxQuestionRunes)))
	}
	if l.svcCtx.AiClient == nil {
		return emit(l.ctx, client, aiprompt.Error("unavailable", "AI 助手暂不可用"))
	}

	sessionID, err := l.ensureSession(identity, req.SessionId, question)
	if err != nil {
		l.Errorf("准备 AI 会话失败: %T", err)
		return emit(l.ctx, client, aiprompt.Error("invalid_session", "会话不可用，请新建会话后重试"))
	}

	if err := emit(l.ctx, client, &types.AiStreamEvent{
		Type:      "meta",
		SessionId: sessionID,
		Guest:     identity.Guest,
		QuotaLeft: aictx.QuotaLeftFrom(l.ctx),
	}); err != nil {
		return err
	}

	history := l.loadHistory(identity, sessionID, req.History)

	if err := emit(l.ctx, client, &types.AiStreamEvent{
		Type: "tool",
		Tool: &types.AiToolEvent{Name: "search_articles", Phase: "started", Label: "检索站内文章"},
	}); err != nil {
		return err
	}

	topK := l.svcCtx.Config.AiConf.Retrieval.TopK
	if topK <= 0 {
		topK = 5
	}
	chunks, err := l.search(question, topK)
	if err != nil {
		l.Errorf("检索文章失败: %T", err)
		chunks = nil
	}

	references, citations := aiprompt.ReferenceBlock(chunks, positiveOr(l.svcCtx.Config.AiConf.Retrieval.MaxCharsTotal, 4000))
	for _, citation := range citations {
		citation.Url = l.citationURL(citation.ArticleId)
	}
	if err := emit(l.ctx, client, &types.AiStreamEvent{
		Type: "tool",
		Tool: &types.AiToolEvent{
			Name:     "search_articles",
			Phase:    "completed",
			Summary:  fmt.Sprintf("命中 %d 篇文章", len(citations)),
			Articles: citations,
		},
	}); err != nil {
		return err
	}

	if references == "" {
		if err := emit(l.ctx, client, &types.AiStreamEvent{Type: "delta", Content: noResultText}); err != nil {
			return err
		}
		if err := emit(l.ctx, client, &types.AiStreamEvent{Type: "done"}); err != nil {
			return err
		}
		// 整轮成功才落库：断开连接时设计上不写会话
		l.persist(sessionID, roleUser, question, 0, 0)
		l.persist(sessionID, roleAssistant, noResultText, 0, 0)
		return nil
	}

	messages := aiprompt.BuildMessages(references, history, question,
		positiveOr(l.svcCtx.Config.AiConf.MaxPromptChars, 12000))

	var answer strings.Builder
	usage, err := l.svcCtx.AiClient.Stream(l.ctx, messages, func(delta string) error {
		answer.WriteString(delta)
		return emit(l.ctx, client, &types.AiStreamEvent{Type: "delta", Content: delta})
	})
	if err != nil {
		event := aiprompt.ErrorEvent(err)
		if event == nil {
			// 客户端断开：上游已随 ctx 取消，不再写会话
			return nil
		}
		var aiErr *ai.Error
		if errors.As(err, &aiErr) {
			l.Errorf("模型生成失败: kind=%s status=%d", aiErr.Kind, aiErr.Status)
		} else {
			l.Errorf("模型生成失败: %T", err)
		}
		return emit(l.ctx, client, event)
	}

	if err := emit(l.ctx, client, &types.AiStreamEvent{
		Type:  "done",
		Usage: &types.AiUsage{PromptTokens: usage.PromptTokens, CompletionTokens: usage.CompletionTokens},
	}); err != nil {
		return err
	}

	// 整轮成功才落库：断开连接时设计上不写会话
	l.persist(sessionID, roleUser, question, 0, 0)
	if text := answer.String(); text != "" {
		l.persist(sessionID, roleAssistant, text, usage.PromptTokens, usage.CompletionTokens)
	}
	return nil
}

// ensureSession 为登录用户建会话；游客不落会话，返回 0。
func (l *AiChatStreamLogic) ensureSession(identity aictx.Identity, requested int64, question string) (int64, error) {
	if identity.Guest {
		return 0, nil
	}
	if requested > 0 {
		// 校验归属：不采信客户端自报的会话 id，也不让它进入 meta
		if _, err := l.svcCtx.AiService.GetSession(l.ctx, &aiservice.GetSessionRequest{Id: requested}); err != nil {
			return 0, err
		}
		return requested, nil
	}

	out, err := l.svcCtx.AiService.CreateSession(l.ctx, &aiservice.CreateSessionRequest{
		Title: aiprompt.TruncateRunes(question, maxTitleRunes),
		Model: l.svcCtx.AiClient.Model(),
	})
	if err != nil {
		return 0, err
	}
	return out.GetSession().GetId(), nil
}

// loadHistory 读取多轮上下文：登录用户取服务端消息，游客取请求里携带的历史。
func (l *AiChatStreamLogic) loadHistory(identity aictx.Identity, sessionID int64, guestHistory []types.AiChatMessage) []ai.Message {
	maxTurns := positiveOr(l.svcCtx.Config.AiConf.MaxHistoryTurns, 6)
	maxChars := positiveOr(l.svcCtx.Config.AiConf.MaxHistoryChars, 4000)

	var history []ai.Message
	if identity.Guest {
		history = capGuestHistory(guestHistory)
	} else if sessionID > 0 {
		out, err := l.svcCtx.AiService.ListMessages(l.ctx, &aiservice.ListMessagesRequest{
			SessionId: sessionID,
			Limit:     int64(maxTurns * 2),
		})
		if err != nil {
			l.Errorf("读取会话历史失败: %T", err)
		} else {
			for _, item := range out.GetList() {
				history = append(history, ai.Message{Role: item.GetRole(), Content: item.GetContent()})
			}
		}
	}

	return aiprompt.BuildHistory(history, maxTurns, maxChars)
}

// capGuestHistory 限制游客携带的历史：最多 50 条、单条 2000 字，非法角色丢弃。
func capGuestHistory(items []types.AiChatMessage) []ai.Message {
	if len(items) > maxGuestHistoryMessages {
		items = items[len(items)-maxGuestHistoryMessages:]
	}

	history := make([]ai.Message, 0, len(items))
	for _, item := range items {
		role := strings.TrimSpace(item.Role)
		content := aiprompt.TruncateRunes(strings.TrimSpace(item.Content), maxGuestMessageRunes)
		if (role != roleUser && role != roleAssistant) || content == "" {
			continue
		}
		history = append(history, ai.Message{Role: role, Content: content})
	}
	return history
}

func (l *AiChatStreamLogic) search(keyword string, topK int) ([]*contentservice.ArticleChunk, error) {
	out, err := l.svcCtx.ContentService.SearchArticles(l.ctx, &contentservice.SearchArticlesRequest{
		Keyword: keyword,
		Limit:   int64(topK),
	})
	if err != nil {
		return nil, err
	}
	return out.GetList(), nil
}

func (l *AiChatStreamLogic) persist(sessionID int64, role, content string, promptTokens, completionTokens int64) {
	if sessionID <= 0 || l.ctx.Err() != nil {
		return
	}
	_, err := l.svcCtx.AiService.AppendMessage(l.ctx, &aiservice.AppendMessageRequest{
		SessionId:        sessionID,
		Role:             role,
		Content:          content,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
	})
	if err != nil {
		l.Errorf("写入会话消息失败: %T", err)
	}
}

func (l *AiChatStreamLogic) citationURL(articleID int64) string {
	base := strings.TrimRight(l.svcCtx.Config.SiteURL, "/")
	return fmt.Sprintf("%s/article/%d", base, articleID)
}

func emit(ctx context.Context, ch chan<- *types.AiStreamEvent, event *types.AiStreamEvent) error {
	select {
	case ch <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func positiveOr(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
