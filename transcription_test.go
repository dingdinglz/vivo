package vivo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestTranscription(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	trans := app.NewTranscription("from.mp3")
	fmt.Println("上传中...")
	e := trans.Upload()
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("开始转写...")
	e = trans.Start()
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("获取任务信息...")
	now, e := trans.GetTaskInfo()
	if e != nil {
		t.Error(e.Error())
		return
	}
	for now != 100 {
		fmt.Println("转写中...进度", now)
		time.Sleep(1 * time.Second)
		now, e = trans.GetTaskInfo()
		if e != nil {
			t.Error(e.Error())
			return
		}
	}

	fmt.Println("转写完成")
	res, e := trans.GetResult()
	if e != nil {
		t.Error(e.Error())
		return
	}
	for _, item := range res {
		fmt.Println("开始秒数：", item.Bg, "结束秒数：", item.Ed, "内容：", item.Onebest)
	}
}
