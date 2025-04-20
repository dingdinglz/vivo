package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestGeoPOI(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	pois, _, e := app.GeoPOISearch("星海广场", "大连市", 1)
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println("检索到：", len(pois), "个")
	for _, item := range pois {
		fmt.Println(item.Name, item.Province+item.City+item.District+item.Address, item.Location, item.Type)
	}
}
