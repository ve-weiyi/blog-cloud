package chatgpt

import (
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func TestChat(t *testing.T) {
	var req []*ChatMessage
	jsonconv.JsonToObject(`[
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
    },
	{
		"role": "assistant",
		"content": "我是人工智能助手，没有具体的名字。你可以随时叫我助手或者小助手。有什么问题可以问我哦？"
   	},
    {
        "role": "user",
        "content": "你可以为自己取一个名字吗？"
    },
	{
		"role": "assistant",
		"content": "当然可以！你可以为我取一个名字吗？"
	},
    {
        "role": "user",
        "content": "我想叫你小助手"
    }
]`, &req)

	res, err := NewAIChatGPT().Chat(req)
	t.Log(err)
	t.Log(jsonconv.ObjectToJsonIndent(res))
}

func TestCosRole(t *testing.T) {
	res, err := NewAIChatGPT().CosRole("担任雅思写作考官")
	t.Log(err)
	t.Log(res)
}

func TestImageGeneration(t *testing.T) {
	res, err := NewAIChatGPT().ImageGeneration("一只猫")
	t.Log(err)
	t.Log(res)
}
