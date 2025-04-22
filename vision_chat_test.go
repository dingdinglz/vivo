package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestVisionChat(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pic, _ := os.ReadFile("test.jpg")
	res, e := app.VisionChat(GenerateRequestID(), GenerateSessionID(), VISION_CHAT_MODEL_BLUELM_PRD, []VisionChatMessage{
		{
			Role:        CHAT_ROLE_USER,
			Content:     GenerateVisionChatImage(pic),
			ContentType: CHAT_MESSAGE_IMAGE,
		},
		{
			Role:        CHAT_ROLE_USER,
			Content:     "描述图片的内容",
			ContentType: CHAT_MESSAGE_TEXT,
		},
	}, nil)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res.Content)
}

func TestVisionChatStream(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pic, _ := os.ReadFile("test.jpg")
	e := app.VisionChatStream(GenerateRequestID(), GenerateSessionID(), VISION_CHAT_MODEL_BLUELM_PRD, []VisionChatMessage{
		{
			Role:        CHAT_ROLE_USER,
			Content:     GenerateVisionChatImage(pic),
			ContentType: CHAT_MESSAGE_IMAGE,
		},
		{
			Role:        CHAT_ROLE_USER,
			Content:     "描述图片的内容",
			ContentType: CHAT_MESSAGE_TEXT,
		},
	}, nil, func(s string) {
		fmt.Print(s)
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
}
