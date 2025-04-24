package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestDialectRecognition(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.DialectRecognition("output.wav")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
