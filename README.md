# vivo

2025中国高校计算机大赛AIGC创新赛，vivo家ai能力第三方SDK

比赛链接：https://aigc.vivo.com.cn/

ai能力文档：https://aigc.vivo.com.cn/#/document/index

DeepWiki（非常好的可以帮你详细理解本仓库的工具）：https://deepwiki.com/dingdinglz/vivo

## 实现能力列表

2025年赛事方提供的全部ai能力已经全部封装完成，使用过程中遇到问题请发issue或者私聊。

如果本包对你有帮助，给个star吧～

- [x] 蓝心大模型70B

- [x] 蓝心大模型多模态

- [x] AI绘画

- [x] 通用OCR

- [x] 实时短语音识别

- [x] 方言自由说

- [x] 同声传译

- [x] 长语音听写

- [x] 长语音转写

- [x] 音频生成

- [x] 声音复制

- [x] 文本翻译

- [x] 文本向量

- [x] 文本相似度

- [x] 查询改写

- [x] 地理编码

## 开始使用

```bash
go get github.com/dingdinglz/vivo
```

## 使用示例

> [!WARNING]
> 下述示例为了保证代码尽可能少，均未做错误处理，实际使用中应进行错误处理！


> [!IMPORTANT]
> 示例中使用的os.Getenv是获取环境变量的意思，为了保证key和id的隐私性，使用时传入id和key的字符串即可，不需要一定要仿照我的os.Getenv


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

#### EasyChat

EasyChat是chat的最简化版本，通过最少的代码量实现效果与Chat相同的操作

你只需要给出问题，就可以获得一个答案，例如：

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
    res, _ := app.EasyChat(vivo.GenerateSessionID(), "你是谁？")
    fmt.Println(res)
}
```

原先通过role=system实现的人设功能可以通过systemPrompt实现，例如

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
    res, _ := app.EasyChat(vivo.GenerateSessionID(), "你是谁？", "你是一只猫娘")
    fmt.Println(res)
}
```

值得一提的是，上下文仍旧生效，具有相同SessionID的所有EasyChat请求共享历史聊天记录，例如：

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
    sessionID := vivo.GenerateSessionID()
    res, _ := app.EasyChat(sessionID, "你是谁？", "你是一只猫娘")
    fmt.Println(res)
    res, _ = app.EasyChat(sessionID, "我刚刚问了你什么？")
    fmt.Println(res)
}
```

这是输出结果：

```
我是猫娘，同时也是你的伙伴和朋友。我可以陪伴你聊天，也可以帮你解决一些问题。
你问我“你是谁？”。
```

可以看出，上下文仍旧有效

> [!NOTE]
> Chat提供Extra参数可以进行更详细的设置，并且在某些需要维护上下文数组的情况下，Chat仍然有效，因此EasyChat并不能直接取代Chat

#### ChatStream && EasyChatStream

以流式形式调用Chat或者EasyChat

示例如下：

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
    app.ChatStream(vivo.GenerateRequestID(), vivo.GenerateSessionID(), []vivo.ChatMessage{
        {
            Role:    vivo.CHAT_ROLE_USER,
            Content: "你是谁？",
        },
    }, nil, func(s string) {
        fmt.Print(s)
    })
}
```

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
    app.EasyChatStream(vivo.GenerateSessionID(), "你是谁？", func(s string) {
        fmt.Print(s)
    })
}
```

### 蓝心大模型多模态

#### VisionChat

相比于Chat，VisionChat可以传入图片，调用参数也仅仅需要把ChatMessage改成VisionMessage即可，需要增加ContentType字段，可以通过vivo.CHAT_MESSAGE_选择，增加了model参数，可以通过vivo.VISION_CHAT_MODEL_选择

``` go
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
	res, _ := app.VisionChat(vivo.GenerateRequestID(), vivo.GenerateSessionID(), vivo.VISION_CHAT_MODEL_BLUELM_PRD, []vivo.VisionChatMessage{
		{
			Role:        vivo.CHAT_ROLE_USER,
			Content:     vivo.GenerateVisionChatImage(pic),
			ContentType: vivo.CHAT_MESSAGE_IMAGE,
		},
		{
			Role:        vivo.CHAT_ROLE_USER,
			Content:     "描述图片的内容",
			ContentType: vivo.CHAT_MESSAGE_TEXT,
		},
	}, nil)
	fmt.Println(res.Content)
}

```

#### VisionChatStream

同理，更改对应的参数即可，形式类似ChatStream

``` go
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
	app.VisionChatStream(vivo.GenerateRequestID(), vivo.GenerateSessionID(), vivo.VISION_CHAT_MODEL_BLUELM_PRD, []vivo.VisionChatMessage{
		{
			Role:        vivo.CHAT_ROLE_USER,
			Content:     vivo.GenerateVisionChatImage(pic),
			ContentType: vivo.CHAT_MESSAGE_IMAGE,
		},
		{
			Role:        vivo.CHAT_ROLE_USER,
			Content:     "描述图片的内容",
			ContentType: vivo.CHAT_MESSAGE_TEXT,
		},
	}, nil, func(s string) {
		fmt.Print(s)
	})
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

### 实时短语音识别 - AsrShortVoiceRecognition

传入文件的地址，注意，文件应当是wav格式，音频格式为16k/16b单声道

给出mp3文件转成该格式的ffmpeg的转换命令

``` bash
ffmpeg -i input.mp3 -ar 16000 -ac 1 -acodec pcm_s16le output.wav
```

返回识别的结果，短语音识别只能识别出一句话，如果一段话的识别，请转向长语音听写或转写

``` go
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
	res, _ := app.AsrShortVoiceRecognition("test.wav")
	fmt.Println(res)
}

```

### 方言自由说 - DialectRecognition

传入文件的地址，注意，文件应当是wav格式，音频格式为16k/16b单声道

给出mp3文件转成该格式的ffmpeg的转换命令

``` bash
ffmpeg -i input.mp3 -ar 16000 -ac 1 -acodec pcm_s16le output.wav
```

返回识别的结果

``` go
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
	res, _ := app.DialectRecognition("test.wav")
	fmt.Println(res)
}

```

### 同声音传译

请使用长语音听写 + 翻译 + 音频生成替代，如果有一定需要流式使用该功能的需要，请发issue。

### 长语音听写 - AsrLongVoiceRecognition

传入文件的地址，注意，文件应当是wav格式，音频格式为16k/16b单声道

给出mp3文件转成该格式的ffmpeg的转换命令

``` bash
ffmpeg -i input.mp3 -ar 16000 -ac 1 -acodec pcm_s16le output.wav
```

返回识别的结果

``` go
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
	res, _ := app.AsrLongVoiceRecognition("test.wav")
	fmt.Println(res)
}

```

### 长语音转写 - Transcription

相较于长语音听写，长语音转写支持更长的音频长度，和更广的音频格式，他支持单次转写文件限制5个小时且小于500M，支持的音频格式有wav，pcm，m4a，mp3，acc，ogg，ogg_opus。

进行长语音转写任务，我们需要经历以下几个步骤：

1. 上传音频
2. 开始转写语音
3. 获取转写进度
4. 获取转写结果

例子写的非常清晰明了，看一遍例子即可明白Transcription的用法

``` go
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
    // 创建转写任务
	trans := app.NewTranscription("output.mp3")
	fmt.Println("开始上传...")
    // 上传音频文件
	e := trans.Upload()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println("启动任务...")
    // 开始转写语音
	e = trans.Start()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	process := 0
	for process != 100 {
		time.Sleep(1 * time.Second)
        // 查询任务进度
		process, e = trans.GetTaskInfo()
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		fmt.Println("当前任务进度：", process, "%")
	}
    // 获取转写结果
	result, e := trans.GetResult()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	for _, item := range result {
		fmt.Println("开始秒数", item.Bg, "结束秒数", item.Ed, "内容：", item.Onebest)
	}
}

```

### 音频生成 - TTS

自动文本转语音（TTS）功能，可将上传的单句文本转成播报音频

第一个参数是生成模式，由vivo.TTS_MODE_选择，第二个参数是音色名称，可以通过声音复刻生成或选择已有音色，在https://aigc.vivo.com.cn/#/document/index?id=1735，查看音色大全，第三个参数是要生成的文本内容。

返回的是pcm格式的数据，可以通过vivo.PcmToWav转成wav格式的数据，用法如下

``` go
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
	res, e := app.TTS(vivo.TTS_MODE_SHORT, "xiaofu", "java是世界上最好的语言")
	if e != nil {
		fmt.Println(e.Error())
	}
	os.WriteFile("test.wav", vivo.PcmToWav(res), os.ModePerm)
}

```

### 声音复制

该服务主要负责将用户上传的录音生成定制的音色（vcn），用户可根据生成的定制音色（vcn）结合短音频生成能力（TTS）合成音频

#### 创建新音色 - VoiceCreate

根据一个音频文件生成一个音色

第一个参数是音频文件的地址，第二个参数是音频文件对应的文字

``` go
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
	vcn, _, _ := app.VoiceCreate("test.wav", "你好，我喜欢玩原神。")
	fmt.Println("vcn:", vcn)
}

```

注意，返回的error!=nil时，不一定是服务调用出现了问题，可能是识别出的文字与传入的文字不一样，可以通过返回的第二个参数进行详细的判断

#### 利用创建的新音色生成一段语音

``` go
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
	vcn, _, _ := app.VoiceCreate("test.wav", "你好，我喜欢玩原神。")
	res, _ := app.TTS(vivo.TTS_MODE_REPLICA, vcn, "你听听这像我嘛？")
	os.WriteFile("test.wav", vivo.PcmToWav(res), os.ModePerm)
}

```

注意mode必须是vivo.TTS_MODE_REPLICA才能正确复刻！

#### 查询一个音色 - VoiceGET

传入的参数是vcn，返回VCNData

``` go
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
	vcn, _ := app.VoiceGET("vcn")
	fmt.Println(vcn.Vcn, vcn.CompleteTime)
}

```

#### 获取已创建的音色列表 - VoiceGetList

``` go
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
	vcnList, _ := app.VoiceGetList()
	for _, item := range vcnList {
		fmt.Println(item.Vcn)
	}
}

```

#### 删除一个音色 - VoiceDelete

传入vcn

``` go
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
	app.VoiceDelete("vcn")
}

```

#### 删除所有音色 - VoiceClean

``` go
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
	app.VoiceClean()
}

```

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

### 文本向量 - TextVector

将用户提供的文本信息表示成计算机可识别的实数向量，用数值向量来表示文本的语义。

第一个参数是使用的模型，由vivo.VECTOR_MODEL_选择，第二个参数是文本列表，返回向量列表

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
    res, _ := app.TextVector(vivo.VECTOR_MODEL_M3E, []string{
        "原神",
        "三国杀",
        "火影忍者",
    })
    fmt.Println(res)
}
```

### 文本相似度 - TextSimilarity

将用户提供的文本信息从语义的角度来判断两者相似度。

第一个参数是使用的模型，由vivo.TEXT_SIMILARITY_MODEL_选择，第二个参数是比对的文本，第三个参数是一系列要比对的句子。

返回结果是对应sentences中每条文本与text文本的相似度

``` go
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
	res, e := app.TextSimilarity(vivo.TEXT_SIMILARITY_MODEL_BGE_LARGE, "科技品牌发展", []string{
		"自动追焦相关报表", "太古汇内云集逾180家知名品牌", "其中逾70个品牌为第一次进驻广州", "交通：商场M层连通地铁三号线石牌桥站；毗邻地铁一号线体育中心站。",
	})
	if e != nil {
		fmt.Println(e.Error())
	}
	fmt.Println(res)
}

```

### 查询改写 - QueryRewrite

查询改写是RAG/AI搜索链路中的重要环节，目的是使用模型对用户当前输入的问题（query）进行理解，并改写为适合搜索引擎检索的query。改写后的结果可根据情况融入历史对话的关键信息，可对复杂问题进行拆解，使得检索召回的知识更加全面、丰富，为最终生成回答提供有力支持。

第一个参数是历史对话记录列表，一个问题，一个回答，最多支持6个，第二个参数是本次询问

返回改写后的查询

``` go
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
	res, e := app.QueryRewrite([]string{
		"战狼2是谁主演的",
		"《战狼2》是由吴京执导并主演的一部军事战争题材电影。影片中，吴京饰演了主角冷锋，他是一名退役的特种部队军人，在非洲执行任务时遭遇了一连串危机和战斗。因此，《战狼2》的主演是吴京。",
	}, "第一部里有他吗")
	if e != nil {
		fmt.Println(e.Error())
	}
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