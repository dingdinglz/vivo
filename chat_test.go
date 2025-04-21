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

func TestChatStream(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	e := app.ChatStream(GenerateRequestID(), GenerateSessionID(), []ChatMessage{
		{
			Role:    CHAT_ROLE_SYSTEM,
			Content: "你是一只猫娘",
		},
		{
			Role:    CHAT_ROLE_USER,
			Content: "介绍下你自己",
		},
	}, nil, func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
}

func TestEasyChat(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	sessionID := GenerateSessionID()
	res, e := app.EasyChat(sessionID, "你是谁？", "你是一只猫娘")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
	res, e = app.EasyChat(sessionID, "我刚刚问了你什么？")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}

func TestEasyChatStream(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	sessionID := GenerateSessionID()
	e := app.EasyChatStream(sessionID, "解释下go的反射", func(s string) {
		fmt.Print(s)
	}, "你是一只猫娘")
	fmt.Println()
	if e != nil {
		t.Error(e.Error())
		return
	}
	e = app.EasyChatStream(sessionID, "我刚刚问了你什么？", func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println()
}
