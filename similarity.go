package vivo

import (
	"encoding/json"
	"errors"
	"io"
)

var (
	TEXT_SIMILARITY_MODEL_BGE_LARGE = "bge-reranker-large"
	TEXT_SIMILARITY_MODEL_BGE_V2    = "bge-reranker-v2-m3"
)

type textSimilarityResponse struct {
	Data []float64 `json:"data"`
}

func (app *Vivo) TextSimilarity(model string, text string, senetences []string) ([]float64, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"model_name": model,
		"query":      text,
		"sentences":  senetences,
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/rerank")
	if e != nil {
		return []float64{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []float64{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []float64{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []float64{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []float64{}, errors.New(msg)
	}
	resData := textSimilarityResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []float64{}, e
	}
	return resData.Data, nil
}
