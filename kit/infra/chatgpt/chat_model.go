package chatgpt

// ChatRequest 是 ChatGPT API 的请求结构体
type ChatRequest struct {
	Model    string         `json:"model"`    // 模型名称
	Messages []*ChatMessage `json:"messages"` // 对话消息列表
}

// ChatMessage 表示对话消息的结构体
type ChatMessage struct {
	Role    string `json:"role"`    // 角色：system 或 user ，assistant ChatGPT 生成的响应
	Content string `json:"content"` // 消息内容
}

// ChatResponse 是 ChatGPT API 的响应结构体
type ChatResponse struct {
	Id      string       `json:"id"`      // 对话 ID
	Object  string       `json:"object"`  // 对象类型
	Created int          `json:"created"` // 创建时间戳
	Model   string       `json:"model"`   // 模型名称
	Choices []ChatChoice `json:"choices"` // 生成的回复列表
	Usage   ChatUsage    `json:"usage"`   // API 调用的使用情况
}

// ChatChoice 表示生成的回复的结构体
type ChatChoice struct {
	Index        int         `json:"index"`         // 回复的索引
	Message      ChatMessage `json:"message"`       // 回复的消息
	FinishReason string      `json:"finish_reason"` // 回复的完成原因
}

// ChatUsage 表示 API 调用的使用情况的结构体
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示 tokens 数量
	CompletionTokens int `json:"completion_tokens"` // 生成回复的 tokens 数量
	TotalTokens      int `json:"total_tokens"`      // 总 tokens 数量
}

// chatgpt角色的结构体
type ChatRole struct {
	Act    string `json:"act"`
	Prompt string `json:"prompt"`
}
