package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestVector(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.TextVector(VECTOR_MODEL_M3E, []string{
		"原神",
		"三国杀",
		"火影忍者",
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
