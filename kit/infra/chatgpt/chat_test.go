package chatgpt

import (
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

func NewGPT() *AIChatGPT {
	return NewAIChatGPT(
		WithApiHost("https://api.openai.com"),
		WithApiKey("xxx"),
		WithModel("gpt-3.5-turbo"),
	)
}

func TestChat(t *testing.T) {
	var req []*ChatMessage
	jsonconv.JsonToAny(`[
    {
        "role": "user",
        "content": "你好"
    },
    {
        "role": "assistant",
        "content": "你好！有什么可以帮助你的吗？如果有任何问题，欢迎问我哦。"
    },
    {
        "role": "user",
        "content": "我想知道你的名字"
    }
]`, &req)

	gpt := NewGPT()
	res, err := gpt.Chat(req)
	t.Log(err)
	t.Log(jsonconv.AnyToJsonIndent(res))
}

func TestCosRole(t *testing.T) {
	gpt := NewGPT()
	res, err := gpt.CosRole("担任雅思写作考官")
	t.Log(err)
	t.Log(res)
}

func TestImageGeneration(t *testing.T) {
	gpt := NewGPT()
	res, err := gpt.ImageGeneration("一只猫")
	t.Log(err)
	t.Log(res)
}
