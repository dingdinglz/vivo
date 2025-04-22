package vivo

import (
	"encoding/json"
	"errors"
	"io"
)

type VoiceCreateResponse struct {
	OpStr     string `json:"op_str"`
	FixOpsStr string `json:"fix_ops_str"`
	ErrorMsg  string `json:"error_msg"`
	ErrorCode int    `json:"error_code"`
	OrgText   string `json:"org_text"`
	AsrText   string `json:"asr_text"`
	VCN       string `json:"vcn"`
}

func (app *Vivo) VoiceCreate(file string, text string) (string, VoiceCreateResponse, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "multipart/form-data")
	httpClient.QueryParams.Add("req_id", "123456")
	httpClient.SetMultipartFormData(map[string]string{
		"text": text,
	})
	httpClient.SetFile("audio", file)
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/replica/create_vcn_task")
	if e != nil {
		return "", VoiceCreateResponse{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return "", VoiceCreateResponse{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return "", VoiceCreateResponse{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return "", VoiceCreateResponse{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return "", VoiceCreateResponse{}, errors.New(msg)
	}
	resData := VoiceCreateResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return "", VoiceCreateResponse{}, e
	}
	if resData.ErrorCode != 0 {
		return "", resData, errors.New(resData.ErrorMsg)
	}
	return resData.VCN, resData, nil
}

type voiceGetResponse struct {
	VcnObj    VCNData `json:"vcn_obj"`
	ErrorMsg  string  `json:"error_msg"`
	ErrorCode int     `json:"error_code"`
}

type VCNData struct {
	CreateTime   string `json:"create_time"`
	CompleteTime string `json:"complete_time"`
	Engineid     string `json:"engineid"`
	TotalCost    int    `json:"total_cost"`
	Process      int    `json:"process"`
	UpdateTime   string `json:"update_time"`
	Vcn          string `json:"vcn"`
}

func (app *Vivo) VoiceGET(vcn string) (VCNData, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.QueryParams.Add("req_id", "123456")
	httpClient.SetBody(map[string]interface{}{
		"vcn": vcn,
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/replica/get_vcn_task")
	if e != nil {
		return VCNData{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return VCNData{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return VCNData{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return VCNData{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return VCNData{}, errors.New(msg)
	}
	resData := voiceGetResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return VCNData{}, e
	}
	if resData.ErrorCode != 0 {
		return VCNData{}, errors.New(resData.ErrorMsg)
	}
	return resData.VcnObj, nil
}

type voiceGetListResponse struct {
	ErrorMsg   string    `json:"error_msg"`
	ErrorCode  int       `json:"error_code"`
	VcnObjList []VCNData `json:"vcn_obj_list"`
}

func (app *Vivo) VoiceGetList() ([]VCNData, error) {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.QueryParams.Add("req_id", "123456")
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/replica/get_vcn_task_list")
	if e != nil {
		return []VCNData{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []VCNData{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []VCNData{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []VCNData{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []VCNData{}, errors.New(msg)
	}
	resData := voiceGetListResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []VCNData{}, e
	}
	if resData.ErrorCode != 0 {
		return []VCNData{}, errors.New(resData.ErrorMsg)
	}
	return resData.VcnObjList, nil
}

type voiceDeleteResponse struct {
	VcnObj    VCNData `json:"vcn_obj"`
	ErrorMsg  string  `json:"error_msg"`
	ErrorCode int     `json:"error_code"`
}

func (app *Vivo) VoiceDelete(vcn string) error {
	httpClient := app.newHttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.QueryParams.Add("req_id", "123456")
	httpClient.SetBody(map[string]interface{}{
		"vcn": vcn,
	})
	httpRes, e := httpClient.Post("https://api-ai.vivo.com.cn/replica/del_task")
	if e != nil {
		return e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return errors.New(msg)
	}
	resData := voiceDeleteResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return e
	}
	if resData.ErrorCode != 0 {
		return errors.New(resData.ErrorMsg)
	}
	return nil
}

func (app *Vivo) VoiceClean() error {
	list, e := app.VoiceGetList()
	if e != nil {
		return e
	}
	for _, item := range list {
		e = app.VoiceDelete(item.Vcn)
		if e != nil {
			return e
		}
	}
	return nil
}
