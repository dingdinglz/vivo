package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestOCR(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pic, _ := os.ReadFile("test.jpg")
	res, e := app.OCR(pic, OCR_MODE_ALL)
	if e != nil {
		t.Error(e.Error())
		return
	}
	switch res.(type) {
	case string:
		fmt.Println(res)
	case []OcrPosData:
		for _, item := range res.([]OcrPosData) {
			fmt.Println(item.Words, item.Location)
		}
	case OcrAllData:
		fmt.Println(res.(OcrAllData))
	default:
		fmt.Println("unkown type")
	}
}
