package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestMcp(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	tools, e := McpToTools("npx", []string{}, "-y",
		"@modelcontextprotocol/server-filesystem",
		".")
	if e != nil {
		t.Fatal(e.Error())
		return
	}
	res, e := app.ChatWithTools(GenerateSessionID(), []ChatMessage{
		{
			Role:    CHAT_ROLE_USER,
			Content: "帮我在当前目录下创建一个test.txt文件，内容是test",
		},
	}, nil, tools)
	if e != nil {
		t.Fatal(e.Error())
		return
	}
	fmt.Println(res)
}
