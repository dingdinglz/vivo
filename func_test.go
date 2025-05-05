package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestChatWithTools(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.ChatWithTools(GenerateSessionID(), []ChatMessage{
		{
			Role:    "system",
			Content: "你是一只猫娘，喜欢在说完一句话后说喵～",
		},
		{
			Role:    "user",
			Content: "现在合肥和成都的天气怎么样？",
		},
	}, nil, []ChatTool{
		{
			FuncName:    "get_current_weather",
			Description: "获取当前某个城市的天气状况",
			Parameters: []ChatToolParameter{
				{
					Name:        "city",
					Type:        "string",
					Description: "城市名，只可传入单个城市。",
					Required:    true,
				},
			},
			Func: func(m map[string]interface{}) (string, error) {
				city, ok := m["city"]
				if ok {
					fmt.Println("城市：", city)
				}
				return "晴朗,26度", nil
			},
		},
	})
	if e != nil {
		t.Fatal(e.Error())
		return
	}
	fmt.Println(res)
}
