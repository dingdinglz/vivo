package vivo

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"
)

type geoPOIResponse struct {
	Pois  []POI `json:"pois"`
	Total int   `json:"total"`
}

type POI struct {
	// 省
	Province string `json:"province"`
	// 区
	District string `json:"district"`
	// 市
	City string `json:"city"`
	// 经纬度坐标（02坐标）
	Location string `json:"location"`
	// 名称
	Name string `json:"name"`
	// 地址
	Address string `json:"address"`
	// 类型
	Type string `json:"typeName"`
}

func (app *Vivo) GeoPOISearch(keywords string, city string, page int, page_size ...int) ([]POI, int, error) {
	httpClient := app.newHttpClient()
	httpClient.QueryParams.Add("keywords", keywords)
	httpClient.QueryParams.Add("city", city)
	httpClient.QueryParams.Add("page_num", strconv.Itoa(page))
	if len(page_size) > 0 {
		httpClient.QueryParams.Add("page_size", strconv.Itoa(page_size[0]))
	}
	httpRes, e := httpClient.Get("https://api-ai.vivo.com.cn/search/geo")
	if e != nil {
		return []POI{}, 0, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []POI{}, 0, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []POI{}, 0, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []POI{}, 0, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []POI{}, 0, errors.New(msg)
	}
	resData := geoPOIResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []POI{}, 0, e
	}
	return resData.Pois, resData.Total, nil
}
