package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestVoiceCreate(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, _, e := app.VoiceCreate("output.wav", "喜欢用雷军AI配音来玩梗恶搞的小朋友们，你们好，我是雷军闲话少说你们要是在用我的声音。来做这些抽象的恶搞视频，我就直接远程操控，你附近的小米酥七撞死你们，谢谢大家。HI。")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
func TestVoiceGet(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	vcn, e := app.VoiceGET("9156908368_761903_3d771654-b6d6-4936-99ed-115228d57fb1_v3")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(vcn)
}

func TestVoiceGetList(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	vcns, e := app.VoiceGetList()
	if e != nil {
		t.Error(e.Error())
		return
	}
	for _, item := range vcns {
		fmt.Println(item.Vcn, item.CompleteTime)
	}
}

func TestVoiceDelete(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	e := app.VoiceDelete("9156908368_761903_4f858468-abee-41e6-8c2d-a584dc07f3cf_v3")
	if e != nil {
		t.Error(e.Error())
		return
	}
}

func TestVoiceClean(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	e := app.VoiceClean()
	if e != nil {
		t.Error(e.Error())
		return
	}
}
