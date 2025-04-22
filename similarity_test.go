package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestTextSimilarity(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.TextSimilarity(TEXT_SIMILARITY_MODEL_BGE_LARGE, "科技品牌发展", []string{
		"自动追焦相关报表", "太古汇内云集逾180家知名品牌", "其中逾70个品牌为第一次进驻广州", "交通：商场M层连通地铁三号线石牌桥站；毗邻地铁一号线体育中心站。",
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
