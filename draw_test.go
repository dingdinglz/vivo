package vivo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestDraw(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	taskID, e := app.Draw("画一只胖猫，在吃麦当劳", DRAW_THEME_GENERAL)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("taskID:", taskID)
	url, code, e := app.DrawGetResult(taskID)
	if e != nil {
		t.Error(e.Error())
		return
	}
	for code == DRAW_TASK_STATUS_QUEUE || code == DRAW_TASK_STATUS_RUNNING {
		time.Sleep(1 * time.Second)
		url, code, e = app.DrawGetResult(taskID)
		if e != nil {
			t.Error(e.Error())
			return
		}
	}
	fmt.Println(url)
}

func TestDraw2Draw(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pic, _ := os.ReadFile("test.jpg")
	taskID, e := app.Draw2Draw(pic, DRAW_THEME_RED, DrawExtra{
		Height: 400,
		Width:  480,
	})
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("taskID:", taskID)
	url, code, e := app.DrawGetResult(taskID)
	if e != nil {
		t.Error(e.Error())
		return
	}
	for code == DRAW_TASK_STATUS_QUEUE || code == DRAW_TASK_STATUS_RUNNING {
		time.Sleep(1 * time.Second)
		url, code, e = app.DrawGetResult(taskID)
		if e != nil {
			t.Error(e.Error())
			return
		}
	}
	fmt.Println(url)
}

func TestDrawGetThemes(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.DrawGetThemes(DRAW_TYPE_TXT2IMG)
	if e != nil {
		t.Error(e.Error())
		return
	}
	for _, item := range res {
		fmt.Println(item.StyleName, item.StyleID)
	}
}

func TestDrawPrompts(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.DrawGetRecommendationPrompts()
	if e != nil {
		t.Error(e.Error())
		return
	}
	for _, item := range res {
		fmt.Println("theme:", item.StyleID)
		fmt.Println("共：", len(item.StylePrompts), "个")
		for _, item2 := range item.StylePrompts {
			fmt.Println(item2.ShortText, item2.LongText)
		}
	}
}

func TestDrawExtend(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pic, _ := os.ReadFile("test.jpg")
	taskID, e := app.DrawExtend(pic, DRAW_THEME_GENERAL, DRAW_EXTEND_MODE_MULTIPLE, DRAW_IMAGE_FORMAT_PNG)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("taskID:", taskID)
	url, code, e := app.DrawGetResult(taskID)
	if e != nil {
		t.Error(e.Error())
		return
	}
	for code == DRAW_TASK_STATUS_QUEUE || code == DRAW_TASK_STATUS_RUNNING {
		time.Sleep(1 * time.Second)
		url, code, e = app.DrawGetResult(taskID)
		if e != nil {
			t.Error(e.Error())
			return
		}
	}
	fmt.Println(url)
}
