package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestAsrShortVoiceRecognition(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.AsrShortVoiceRecognition("output.wav")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}

func TestAsrLongVoiceRecognition(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.AsrLongVoiceRecognition("output.wav")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
