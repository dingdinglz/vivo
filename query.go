package vivo

import (
	"encoding/json"
	"errors"
	"io"
)

type queryRewriteResponse struct {
	Result []string `json:"result"`
	Code   int      `json:"code"`
}

func (app *Vivo) QueryRewrite(historys []string, query string) (string, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	for len(historys) < 6 {
		historys = append([]string{""}, historys...)
	}
	if len(historys) > 6 {
		historys = historys[len(historys)-6:]
	}
	httpClient.SetBody(map[string]interface{}{
		"prompts": [][]string{
			historys,
			{
				query,
			},
		},
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/query_rewrite_base")
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
	resData := queryRewriteResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return "", e
	}
	if resData.Code != 0 {
		errMap := make(map[int]string)
		errMap[-2] = "请求列表格式错误"
		errMap[-3] = "当前query长度大于50"
		errMap[-4] = "当前query含有特定词语（A类）"
		errMap[-5] = "当前query含有特定词语（B类）"
		errMap[-6] = "上轮历史只有query或只有answer"
		errMap[-8] = "当前query含有特定模版不进行改写"
		errMap[-9] = "模型判定无需改写"
		errMap[-3002] = "服务运行异常"
		return "", errors.New(errMap[resData.Code])
	}
	if len(resData.Result) == 0 {
		return "", errors.New("unkown error")
	}
	return resData.Result[0], nil
}
