package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestChat(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	sessionID := GenerateSessionID()
	messages := []ChatMessage{
		{
			Role:    CHAT_ROLE_USER,
			Content: "介绍一下你自己",
		}}
	resMessage, e := app.Chat(GenerateRequestID(), sessionID, messages, nil)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(resMessage)
	messages = append(messages, resMessage, ChatMessage{
		Role:    CHAT_ROLE_USER,
		Content: "我刚刚问了你什么？",
	})
	resMessage, e = app.Chat(GenerateRequestID(), sessionID, messages, nil)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(resMessage)
}
