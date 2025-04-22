package vivo

import (
	"fmt"
	"os"
	"testing"
)

func TestQueryRewrite(t *testing.T) {
	app := NewVivoAIGC(Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, e := app.QueryRewrite([]string{
		"战狼2是谁主演的",
		"《战狼2》是由吴京执导并主演的一部军事战争题材电影。影片中，吴京饰演了主角冷锋，他是一名退役的特种部队军人，在非洲执行任务时遭遇了一连串危机和战斗。因此，《战狼2》的主演是吴京。",
	}, "第一部里有他吗")
	if e != nil {
		t.Error(e.Error())
		return
	}
	fmt.Println(res)
}
