package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.Translate(TRANSLATE_LANGUAGE_CHINESE, TRANSLATE_LANGUAGE_ENGLISH, "你好，我爱吃香菜")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
