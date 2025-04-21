package vivo

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	VECTOR_MODEL_M3E = "m3e-base"
)

type textVectorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    [][]float64 `json:"data"`
}

func (app *Vivo) TextVector(model string, text []string) ([][]float64, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"model_name": model,
		"sentences":  text,
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/embedding-model-api/predict/batch")
	if e != nil {
		return [][]float64{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return [][]float64{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return [][]float64{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return [][]float64{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return [][]float64{}, errors.New(msg)
	}
	resData := textVectorResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return [][]float64{}, e
	}
	if resData.Code != 0 {
		return [][]float64{}, errors.New(resData.Message)
	}
	return resData.Data, nil
}
