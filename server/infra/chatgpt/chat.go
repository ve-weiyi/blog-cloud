package chatgpt

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/httpx"
)

const (
	RoleUser = "user"
	RoleAI   = "assistant"

	ChatUrl  = "/v1/chat/completions"
	ImageUrl = "/v1/images/generations"
)

type ChatGPT interface {
	Chat(req []*ChatMessage) (resp *ChatResponse, err error)
}

type AIChatGPT struct {
	ApiHost string
	ApiKey  string
	Model   string
}

func NewAIChatGPT() *AIChatGPT {
	return &AIChatGPT{
		ApiHost: global.CONFIG.ChatGPT.ApiHost,
		ApiKey:  global.CONFIG.ChatGPT.ApiKey,
		Model:   global.CONFIG.ChatGPT.Model,
	}
}
func (s *AIChatGPT) Chat(req []*ChatMessage) (resp *ChatResponse, err error) {
	content := ChatRequest{
		Model:    s.Model,
		Messages: req,
	}

	fmt.Println("content", content)

	res, err := httpx.NewClient(
		httpx.WithHeader("Content-Type", "application/json"),
		httpx.WithHeader("Authorization", "Bearer "+s.ApiKey),
		httpx.WithBodyObject(content),
	).DoRequest("POST", s.ApiHost+ChatUrl)
	if err != nil {
		return nil, err
	}

	resp = &ChatResponse{}
	err = json.Unmarshal([]byte(res), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AIChatGPT) ImageGeneration(prompt string) (resp any, err error) {

	req := map[string]any{
		"model":  "dall-e-3",
		"prompt": prompt,
		"size":   "1024x1024",
		"n":      1,
	}

	res, err := httpx.NewClient(
		httpx.WithHeader("Content-Type", "application/json"),
		httpx.WithHeader("Authorization", "Bearer "+s.ApiKey),
		httpx.WithBodyObject(req),
	).DoRequest("POST", s.ApiHost+ImageUrl)
	if err != nil {
		return nil, err
	}

	//resp = &ImageResponse{}
	//err = json.Unmarshal([]byte(res), resp)
	//if err != nil {
	//	return nil, err
	//}

	return string(res), nil
}

func (s *AIChatGPT) CosRole(act string) (resp *ChatResponse, err error) {

	prompt, err := s.getRole(act)
	if err != nil {
		return nil, err
	}

	content := ChatRequest{
		Model: s.Model,
		Messages: append([]*ChatMessage{{
			Role:    RoleUser,
			Content: prompt,
		}}),
	}

	res, err := httpx.NewClient(
		httpx.WithHeader("Content-Type", "application/json"),
		httpx.WithHeader("Authorization", "Bearer "+s.ApiKey),
		httpx.WithBodyObject(content),
	).DoRequest("POST", s.ApiHost+ChatUrl)
	if err != nil {
		return nil, err
	}

	resp = &ChatResponse{}
	err = json.Unmarshal([]byte(res), resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AIChatGPT) getRole(act string) (string, error) {
	roles, err := s.readModelEmbed()
	//roles, err := s.readModelJSON("./prompts-zh.json")
	if err != nil {
		return "", err
	}

	for _, role := range roles {
		if role.Act == act {
			return role.Prompt, nil
		}
	}

	return "", fmt.Errorf("not found role")
}

//go:embed prompts-zh.json
var rolePrompts string

func (s *AIChatGPT) readModelEmbed() (roles []*ChatRole, err error) {
	// 解析 JSON 数据到结构体
	err = json.Unmarshal([]byte(rolePrompts), &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return roles, nil
}

func (s *AIChatGPT) readModelJSON(filepath string) (roles []*ChatRole, err error) {
	//读取 JSON 文件内容
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %v", err)
	}

	// 解析 JSON 数据到结构体
	err = json.Unmarshal(jsonData, &roles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return roles, nil
}
