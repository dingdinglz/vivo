package vivo

import (
	"os"
	"testing"
)

func TestTTS(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pcmData, e := app.TTS(TTS_MODE_LONG, "x2_yige", "go语言是世界上最好的语言")
	if e != nil {
		t.Error(e.Error())
		return
	}
	os.WriteFile("test.wav", PcmToWav(pcmData), os.ModePerm)
}
