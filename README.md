# vivo

2025中国高校计算机大赛AIGC创新赛，vivo家ai能力第三方SDK

比赛链接：https://aigc.vivo.com.cn/

ai能力文档：https://aigc.vivo.com.cn/#/document/index

## 实现能力列表

- [x] 蓝心大模型70B

- [ ] 蓝心大模型多模态

- [x] AI绘画

- [x] 通用OCR

- [ ] 实时短语音识别

- [ ] 方言自由说

- [ ] 同声传译

- [ ] 长语音听写

- [ ] 长语音转写

- [ ] 音频生成

- [ ] 声音复制

- [x] 文本翻译

- [ ] 文本向量

- [ ] 文本相似度

- [ ] 查询改写

- [x] 地理编码

## 开始使用

```bash
go get github.com/dingdinglz/vivo
```

## 使用示例

> [!WARNING]
> 下述示例为了保证代码尽可能少，均未做错误处理，实际使用中应进行错误处理！

### 蓝心大模型70B

#### Chat

标准对话

Chat的第一个参数是Request_id，每次请求用GenerateRequestID生成一个即可，SessionID可以用GenerateSessionID生成，每次对话需要保持一致，比如上下文，Message是一个数组，Role由vivo.CHAT_ROLE_表示四种身份，Content是对话内容

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    res, _ := app.Chat(vivo.GenerateRequestID(), vivo.GenerateSessionID(), []vivo.ChatMessage{
        {
            Role:    vivo.CHAT_ROLE_USER,
            Content: "介绍下你自己吧",
        },
    }, nil)
    fmt.Println(res.Content)
} 
```

上面是无上下文示例，下面是有上下文示例

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    session_id := vivo.GenerateSessionID()
    messages := []vivo.ChatMessage{
        {
            Role:    vivo.CHAT_ROLE_USER,
            Content: "介绍下你自己吧",
        },
    }
    res, _ := app.Chat(vivo.GenerateRequestID(), session_id, messages, nil)
    fmt.Println(res.Content)
    messages = append(messages, res)
    messages = append(messages, vivo.ChatMessage{
        Role:    vivo.CHAT_ROLE_USER,
        Content: "我刚刚问了你什么？",
    })
    res, _ = app.Chat(vivo.GenerateRequestID(), session_id, messages, nil)
    fmt.Println(res.Content)
}
```

### AI绘画

#### 文生图 - Draw

第一个参数是提示词，用于描述要画的内容，第二个参数是要使用的绘画风格，你可以使用vivo.DRAW_THEME_选择，也可以用DrawGetThemes获取，第三个参数是可选参数，你可以用它来控制一些额外参数，例如大小，负面提示词等等

画图结束后返回taskid，需要通过DrawGetResult获取当前绘图状态

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    task_id, _ := app.Draw("画一只猫猫", vivo.DRAW_THEME_GENERAL)
    url, status, _ := app.DrawGetResult(task_id)
    for status == vivo.DRAW_TASK_STATUS_QUEUE || status == vivo.DRAW_TASK_STATUS_RUNNING {
        time.Sleep(1 * time.Second)
        url, status, _ = app.DrawGetResult(task_id)
    }
    fmt.Println(url)
}
```

#### 图生图 - Draw2Draw

第一个参数是图片，用于描述要画的内容，第二个参数是要使用的绘画风格，你可以使用vivo.DRAW_THEME_选择，也可以用DrawGetThemes获取，第三个参数是可选参数，你可以用它来控制一些额外参数，例如大小，负面提示词等等

画图结束后返回taskid，需要通过DrawGetResult获取当前绘图状态

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    pic, _ := os.ReadFile("test.png")
    task_id, _ := app.Draw2Draw(pic, vivo.DRAW_THEME_GENERAL)
    url, status, _ := app.DrawGetResult(task_id)
    for status == vivo.DRAW_TASK_STATUS_QUEUE || status == vivo.DRAW_TASK_STATUS_RUNNING {
        time.Sleep(1 * time.Second)
        url, status, _ = app.DrawGetResult(task_id)
    }
    fmt.Println(url)
}
```

#### 取消绘画任务 - DrawCancel

你可以在绘画未完成之前取消该任务

```go
package main

import (
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    app.DrawCancel("task_id")
}
```

#### 获取风格列表 - DrawGetThemes

你可以通过绘画的类型获取所有风格，传入绘画的类型即可，由vivo.DRAW_TYPE_获得

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    themes, _ := app.DrawGetThemes(vivo.DRAW_TYPE_IMG2IMG)
    for _, item := range themes {
        fmt.Println(item.StyleName, item.StyleID)
    }
}
```

#### 图片外拓 - DrawExtend

让ai进行图片拓展，第一个参数是图片，第二个参数是绘画风格，第三个参数是拓展模式，通过vivo.DRAW_EXTEND_MODE_选择，第四个参数是生成的图片格式，由vivo.DRAW_IMAGE_FORMAT_选择，第五个是可选参数，用于控制额外的参数，例如图片的大小等等

返回一个task_id，通过DrawGetResult获取绘画状态

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    pic, _ := os.ReadFile("test.png")
    task_id, _ := app.DrawExtend(pic, "dd72470d24fd59bc9df42046c4b27bae", vivo.DRAW_EXTEND_MODE_MULTIPLE, vivo.DRAW_IMAGE_FORMAT_PNG)
    url, status, _ := app.DrawGetResult(task_id)
    for status == vivo.DRAW_TASK_STATUS_QUEUE || status == vivo.DRAW_TASK_STATUS_RUNNING {
        time.Sleep(1 * time.Second)
        url, status, _ = app.DrawGetResult(task_id)
    }
    fmt.Println(url)
}
```

### 通用ocr - OCR

识别用户向服务请求的某张图中的所有文字，并返回文字在图片中的位置信息，方便用户进行文字排版的二次处理参考。

第一个参数是图片，第二个参数是识别模式，通过vivo.OCR_MODE_选择，共支持三种模式，返回的格式均不相同，具体见示例

#### 仅返回文字信息

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    pic, _ := os.ReadFile("test.png")
    res, _ := app.OCR(pic, vivo.OCR_MODE_ONLY)
    fmt.Println(res.(string))
}
```

该模式下返回的类型是string

#### 提供文字信息和坐标信息

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    pic, _ := os.ReadFile("test.png")
    res, _ := app.OCR(pic, vivo.OCR_MODE_POS)
    for _, item := range res.([]vivo.OcrPosData) {
        fmt.Println(item.Words, item.Location)
    }
}
```

该模式下返回的是[]vivo.OcrPosData

#### 混合提供

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    pic, _ := os.ReadFile("test.png")
    res, _ := app.OCR(pic, vivo.OCR_MODE_ALL)
    fmt.Println(res.(vivo.OcrAllData))
}
```

该模式下返回的是vivo.OcrAllData，具有两个参数，分别对应上面两种模式的返回参数

### 文本翻译 - Translate

将一段源语言文本转换成目标语言文本，可根据语言参数的不同实现多国语言之间的互译。

第一个参数是传入文本的语言，第二个参数是翻译到的文本的语言，这两个参数可以通过vivo.TRANSLATE_LANGUAGE_指定，第三个参数是传入的文本，返回翻译后的文本

```go
package main

import (
	"fmt"
	"os"

	"github.com/dingdinglz/vivo"
)

func main() {
	app := vivo.NewVivoAIGC(vivo.Config{
		AppID:  os.Getenv("APPID"),
		AppKey: os.Getenv("APPKEY"),
	})
	res, _ := app.Translate(vivo.TRANSLATE_LANGUAGE_CHINESE, vivo.TRANSLATE_LANGUAGE_JAPANESE, "我不吃香菜")
	fmt.Println(res)
}

```



### 地理编码

#### POI搜索 - GeoPOISearch

输入关键字，查询对应城市的POI接口，输出相关联的地理名称、类别、经度纬度、附近的酒店饭店商铺等信息。

第一个参数是地点，第二个参数是城市名，第三个参数是当前显示的页数，第四个是可选参数，一页有多少个地点

```go
package main

import (
    "fmt"
    "os"

    "github.com/dingdinglz/vivo"
)

func main() {
    app := vivo.NewVivoAIGC(vivo.Config{
        AppID:  os.Getenv("APPID"),
        AppKey: os.Getenv("APPKEY"),
    })
    _, total, _ := app.GeoPOISearch("星海广场", "大连市", 1, 10)
    fmt.Println("共检索到", total, "个地点")
    pages := int(total / 10)
    if total%pages != 0 {
        pages++
    }
    for i := 1; i <= pages; i++ {
        poi, _, _ := app.GeoPOISearch("星海广场", "大连市", i, 10)
        for _, item := range poi {
            fmt.Println(item.Name, item.Province+item.City+item.District+item.Address, item.Location, item.Type)
        }
    }
}
```