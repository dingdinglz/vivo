package vivo

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	TRANSLATE_LANGUAGE_AUTO     = "auto" // （我测试的时候似乎不能用，还是手动选择语言吧）自动识别语言
	TRANSLATE_LANGUAGE_CHINESE  = "zh-CHS"
	TRANSLATE_LANGUAGE_ENGLISH  = "en"
	TRANSLATE_LANGUAGE_JAPANESE = "ja"
	TRANSLATE_LANGUAGE_KOREAN   = "ko" // 韩语
)

type translateResponse struct {
	Code int `json:"code"`
	Data struct {
		Translation string `json:"translation"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func (app *Vivo) Translate(from string, to string, text string) (string, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpClient.SetFormData(map[string]string{
		"from":      from,
		"to":        to,
		"text":      text,
		"app":       "test",
		"requestId": GenerateRequestID(),
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/translation/query/self")
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
	resData := translateResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return "", e
	}
	if resData.Code != 0 {
		return "", errors.New(resData.Msg)
	}
	return resData.Data.Translation, nil
}
