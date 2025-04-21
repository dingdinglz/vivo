package vivo

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

var (
	OCR_MODE_ONLY = 0 // 仅返回文字信息
	OCR_MODE_POS  = 1 // 提供文字信息和坐标信息
	OCR_MODE_ALL  = 2 // ONLY和POS两种模式的信息同时提供
)

type ocrResponse1 struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Result    struct {
		Words []struct {
			Words string `json:"words"`
		} `json:"words"`
	} `json:"result"`
}

type ocrResponse2 struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Result    struct {
		Ocr []OcrPosData `json:"OCR"`
	} `json:"result"`
}

type OcrPosData struct {
	Location OcrPosDataLocation `json:"location"`
	Words    string             `json:"words"`
}

type OcrPosDataLocation struct {
	DownLeft  OcrPoint `json:"down_left"`
	DownRight OcrPoint `json:"down_right"`
	TopLeft   OcrPoint `json:"top_left"`
	TopRight  OcrPoint `json:"top_right"`
}

type OcrPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ocrResponse3 struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Result    struct {
		Ocr   []OcrPosData `json:"OCR"`
		Words []struct {
			Words string `json:"words"`
		} `json:"words"`
	} `json:"result"`
}

type OcrAllData struct {
	Pos  []OcrPosData
	Word string
}

func (app *Vivo) OCR(pic []byte, mode int) (interface{}, error) {
	if mode < 0 || mode > 2 {
		return nil, errors.New("unsupported mode")
	}
	httpClient := app.newHttpClient()
	httpClient.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	httpClient.SetFormData(map[string]string{
		"image":      base64encode(pic),
		"pos":        strconv.Itoa(mode),
		"businessid": "1990173156ceb8a09eee80c293135279",
	})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/ocr/general_recognition")
	if e != nil {
		return nil, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return nil, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return nil, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return nil, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return nil, errors.New(msg)
	}
	switch mode {
	case OCR_MODE_ONLY:
		resData := ocrResponse1{}
		e := json.Unmarshal(body, &resData)
		if e != nil {
			return nil, e
		}
		if resData.ErrorCode != 0 {
			return nil, errors.New(resData.ErrorMsg)
		}
		res := ""
		for _, item := range resData.Result.Words {
			res += item.Words + "\n"
		}
		return res, nil
	case OCR_MODE_POS:
		resData := ocrResponse2{}
		e := json.Unmarshal(body, &resData)
		if e != nil {
			return nil, e
		}
		if resData.ErrorCode != 0 {
			return nil, errors.New(resData.ErrorMsg)
		}
		return resData.Result.Ocr, nil
	case OCR_MODE_ALL:
		resData := ocrResponse3{}
		e := json.Unmarshal(body, &resData)
		if e != nil {
			return nil, e
		}
		if resData.ErrorCode != 0 {
			return nil, errors.New(resData.ErrorMsg)
		}
		res := ""
		for _, item := range resData.Result.Words {
			res += item.Words + "\n"
		}
		return OcrAllData{
			Word: res,
			Pos:  resData.Result.Ocr,
		}, nil
	default:
		return nil, errors.New("unsupported mode")
	}
}
