package vivo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	DRAW_THEME_GENERAL = "4cbc9165bc615ea0815301116e7925a3" // 通用v6.0
	DRAW_THEME_ANIME   = "85ae2641576f5c409b273e0f490f15c0" // 梦幻动漫
	DRAW_THEME_REAL    = "85062a504de85d719df43f268199c308" // 唯美写实
	DRAW_THEME_RED     = "b3aacd62d38c5dbfb3f3491c00ba62f0" // 绯红烈焰
	DRAW_THEME_JAPEN   = "897c280803be513fa947f914508f3134" // 彩绘日漫

	DRAW_TASK_STATUS_QUEUE    = 0 // 队列中，等待处理
	DRAW_TASK_STATUS_RUNNING  = 1 // 正在处理
	DRAW_TASK_STATUS_COMPLETE = 2 // 处理完成
	DRAW_TASK_STATUS_FAILED   = 3 // 处理失败
	DRAW_TASK_STATUS_CANCEL   = 4 // 已取消

	DRAW_TYPE_TXT2IMG = "txt2img" // 文生图
	DRAW_TYPE_IMG2IMG = "img2img" // 图生图

	DRAW_EXTEND_MODE_MULTIPLE   = 1 // 外扩倍数扩充
	DRAW_EXTEND_MODE_PROPORTION = 2 // 按比例扩充
	DRAW_EXTEND_MODE_PIXEL      = 3 // 按像素扩充

	DRAW_IMAGE_FORMAT_PNG  = "PNG"
	DRAW_IMAGE_FORMAT_JPEG = "JPEG"
)

type drawRequest struct {
	UserAccount string `json:"userAccount"`
	Prompt      string `json:"prompt,omitempty"`
	InitImages  string `json:"initImages,omitempty"`
	ImageType   int    `json:"imageType"`
	StyleConfig string `json:"styleConfig,omitempty"`
	DrawExtra
}

type DrawExtra struct {
	Height            int     `json:"height,omitempty"`
	Width             int     `json:"width,omitempty"`
	Seed              int     `json:"seed,omitempty"`
	CfgScale          float64 `json:"cfgScale,omitempty"`
	DenoisingStrength float64 `json:"denoisingStrength,omitempty"`
	CtrlNetStrength   float64 `json:"ctrlNetStrength,omitempty"`
	Steps             int     `json:"steps,omitempty"`
	NegativePrompt    string  `json:"negativePrompt,omitempty"`
}

type drawResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		TaskID string `json:"task_id"`
	} `json:"result"`
}

// Draw 文生图，返回任务id，theme可通过内置的DRAW_THEME_获取，或者通过DrawGetThemes获取
func (app *Vivo) Draw(prompt string, theme string, extra ...DrawExtra) (string, error) {
	client := app.newHttpClient()
	client.Header.Set("Content-Type", "application/json")
	requestBody := drawRequest{}
	if len(extra) != 0 {
		requestBody.DrawExtra = extra[0]
	}
	requestBody.UserAccount = "vivoaigccontest"
	requestBody.Prompt = prompt
	requestBody.StyleConfig = theme
	client.SetBody(requestBody)
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/api/v1/task_submit")
	if e != nil {
		return "", e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return "", e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return "", errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return "", errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return "", errors.New(msg)
	}
	drawRes := drawResponse{}
	e = json.Unmarshal(body, &drawRes)
	if e != nil {
		return "", e
	}
	if drawRes.Code != 200 {
		return "", errors.New(drawRes.Msg)
	}
	return drawRes.Result.TaskID, nil
}

type drawResultResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		ImagesURL []string `json:"images_url"`
		Status    int      `json:"status"`
	} `json:"result"`
}

// DrawGetResult 如果成功则返回url，其他状态（可能是等待、排队等等）返回string为空，返回的int是状态
func (app *Vivo) DrawGetResult(task_id string) (string, int, error) {
	httpClient := app.newHttpClient()
	httpClient.QueryParams.Add("task_id", task_id)
	httpRes, e := httpClient.Get("https://api-ai.vivo.com.cn/api/v1/task_progress")
	if e != nil {
		return "", 0, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return "", 0, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return "", 0, e
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return "", 0, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return "", 0, errors.New(msg)
	}
	resData := drawResultResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return "", 0, e
	}
	imageUrl := ""
	if len(resData.Result.ImagesURL) > 0 {
		imageUrl = resData.Result.ImagesURL[0]
	}
	return imageUrl, resData.Result.Status, nil
}

// Draw2Draw 图生图 pic是要重画的文件
func (app *Vivo) Draw2Draw(pic []byte, theme string, extra ...DrawExtra) (string, error) {
	client := app.newHttpClient()
	client.Header.Set("Content-Type", "application/json")
	requestBody := drawRequest{}
	requestBody.UserAccount = "vivoaigccontest"
	if len(extra) != 0 {
		requestBody.DrawExtra = extra[0]
	}
	requestBody.InitImages = "data:image/png;base64," + base64encode(pic)
	requestBody.StyleConfig = theme
	client.SetBody(requestBody)
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/api/v1/task_submit")
	if e != nil {
		return "", e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return "", e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return "", errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return "", errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return "", errors.New(msg)
	}
	drawRes := drawResponse{}
	e = json.Unmarshal(body, &drawRes)
	if e != nil {
		return "", e
	}
	if drawRes.Code != 200 {
		return "", errors.New(drawRes.Msg)
	}
	return drawRes.Result.TaskID, nil
}

func (app *Vivo) DrawCancel(task_id string) {
	httpClient := app.newHttpClient()
	httpClient.Header.Add("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"task_id": task_id,
	})
	httpClient.Post("https://api-ai.vivo.com.cn/api/v1/task_cancel")
}

type drawThemesResponse struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Result []DrawTheme `json:"result"`
}

type DrawTheme struct {
	StyleID   string `json:"style_id"`
	StyleName string `json:"style_name"`
}

// DrawGetThemes获取绘图全部theme ， StyleID对应theme参数，DRAW_THEME_内置了一些theme
func (app *Vivo) DrawGetThemes(drawType string) ([]DrawTheme, error) {
	httpClient := app.newHttpClient()
	httpClient.QueryParams.Add("styleType", drawType)
	httpRes, e := httpClient.Get("http://api-ai.vivo.com.cn/api/v1/styles")
	if e != nil {
		return []DrawTheme{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []DrawTheme{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []DrawTheme{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []DrawTheme{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []DrawTheme{}, errors.New(msg)
	}
	resData := drawThemesResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []DrawTheme{}, e
	}
	if resData.Code != 200 {
		return []DrawTheme{}, errors.New(resData.Msg)
	}
	return resData.Result, nil
}

type drawPromptsResponse struct {
	Code   int                `json:"code"`
	Msg    string             `json:"msg"`
	Result []DrawThemePrompts `json:"result"`
}

type DrawThemePrompts struct {
	StylePrompts []DrawPrompt `json:"style_prompts"`
	StyleID      string       `json:"style_id"`
}

type DrawPrompt struct {
	LongText  string `json:"long_text"`
	ShortText string `json:"short_text"`
}

// DrawGetRecommendationPrompts 获取文生图推荐词列表
func (app *Vivo) DrawGetRecommendationPrompts() ([]DrawThemePrompts, error) {
	httpClient := app.newHttpClient()
	httpRes, e := httpClient.Get("https://api-ai.vivo.com.cn/api/v1/prompts")
	if e != nil {
		return []DrawThemePrompts{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []DrawThemePrompts{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []DrawThemePrompts{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []DrawThemePrompts{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []DrawThemePrompts{}, errors.New(msg)
	}
	resData := drawPromptsResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []DrawThemePrompts{}, e
	}
	if resData.Code != 200 {
		return []DrawThemePrompts{}, errors.New(resData.Msg)
	}
	fmt.Println(string(body))
	return resData.Result, nil
}

type drawExtendRequest struct {
	InitImages   string `json:"initImages"`
	ImageType    int    `json:"imageType"`
	StyleConfig  string `json:"styleConfig"`
	OutpaintMode int    `json:"outpaintMode"`
	ImageFormat  string `json:"imageFormat"`
	DrawExtendExtra
}

type DrawExtendExtra struct {
	Seed      int                       `json:"seed,omitempty"`
	PadFactor int                       `json:"padFactor,omitempty"`
	PadRatio  string                    `json:"padRatio,omitempty"`
	PadPixel  *DrawExtendPadPixelConfig `json:"padPixel,omitempty"`
}

type DrawExtendPadPixelConfig struct {
	PadUp    int `json:"pad_up,omitempty"`
	PadDown  int `json:"pad_down,omitempty"`
	PadLeft  int `json:"pad_left,omitempty"`
	PadRight int `json:"pad_right,omitempty"`
	PadAll   int `json:"pad_all,omitempty"`
}

func (app *Vivo) DrawExtend(pic []byte, theme string, mode int, format string, extra ...DrawExtendExtra) (string, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	requestBody := drawExtendRequest{
		InitImages:   "data:image/png;base64," + base64encode(pic),
		StyleConfig:  "dd72470d24fd59bc9df42046c4b27bae",
		OutpaintMode: mode,
		ImageFormat:  format,
	}
	httpClient.SetBody(requestBody)
	if len(extra) > 0 {
		requestBody.DrawExtendExtra = extra[0]
	}
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/api/v1/outpaint_task_submit")
	if e != nil {
		return "", e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return "", e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return "", errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return "", errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return "", errors.New(msg)
	}
	drawRes := drawResponse{}
	e = json.Unmarshal(body, &drawRes)
	if e != nil {
		return "", e
	}
	if drawRes.Code != 200 {
		return "", errors.New(drawRes.Msg)
	}
	return drawRes.Result.TaskID, nil
}
