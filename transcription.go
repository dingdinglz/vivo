package vivo

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"resty.dev/v3"
)

var (
	transcription_chunk_size int64 = 5 * 1024 * 1024
)

type Transcription struct {
	sessionID        string
	taskID           string
	filePath         string
	createhttpClient func() *resty.Request
	slice_num        int
	audio_id         string
	file_size        int64
}

func (app *Vivo) NewTranscription(filePath string) *Transcription {
	return &Transcription{
		sessionID: GenerateSessionID(),
		filePath:  filePath,
		createhttpClient: func() *resty.Request {
			httpClient := app.newHttpClient()
			addPublicTranscriptionQuery(httpClient)
			return httpClient
		},
	}
}

func addPublicTranscriptionQuery(httpClient *resty.Request) {
	httpClient.QueryParams.Add("client_version", "2.0")
	httpClient.QueryParams.Add("package", "pack")
	httpClient.QueryParams.Add("user_id", "2addc42b7ae689dfdf1c63e220df52a2")
	httpClient.QueryParams.Add("system_time", int64toString(time.Now().Unix()))
	httpClient.QueryParams.Add("engineid", "fileasrrecorder")
}

type createVideoResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Data struct {
		AudioID string `json:"audio_id"`
	} `json:"data"`
}

func (t *Transcription) createVideo() error {
	fileInfo, e := os.Stat(t.filePath)
	if e != nil {
		return e
	}
	t.file_size = fileInfo.Size()
	t.slice_num = int(fileInfo.Size() / transcription_chunk_size)
	if fileInfo.Size()%transcription_chunk_size != 0 {
		t.slice_num++
	}
	httpClient := t.createhttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"audio_type":  "auto",
		"x-sessionId": t.sessionID,
		"slice_num":   t.slice_num,
	})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/lasr/create")
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
	resData := createVideoResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return e
	}
	if resData.Code != 0 {
		return errors.New(resData.Desc)
	}
	t.audio_id = resData.Data.AudioID
	return nil
}

type uploadChunkResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (t *Transcription) uploadChunk(fileBody []byte, index int) error {
	httpClient := t.createhttpClient()
	httpClient.Header.Set("Content-Type", "multipart/form-data")
	httpClient.QueryParams.Add("audio_id", t.audio_id)
	httpClient.QueryParams.Add("x-sessionId", t.sessionID)
	httpClient.QueryParams.Add("slice_index", strconv.Itoa(index))
	httpClient.SetMultipartFormData(map[string]string{"file": string(fileBody)})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/lasr/upload")
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
	resData := uploadChunkResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return e
	}
	if resData.Code != 0 {
		return errors.New(resData.Desc)
	}
	return nil
}

func (t *Transcription) Upload() error {
	t.createVideo()
	file, e := os.Open(t.filePath)
	if e != nil {
		return e
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	i := 0
	for {
		chunk := make([]byte, transcription_chunk_size)
		n, e := reader.Read(chunk)
		if e != nil {
			if e == io.EOF {
				break
			}
			return e
		}
		chunk = chunk[:n]
		e = t.uploadChunk(chunk, i)
		if e != nil {
			return e
		}
		i++
	}
	return nil
}

type startResponse struct {
	Data struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (t *Transcription) Start() error {
	httpClient := t.createhttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"audio_id":    t.audio_id,
		"x-sessionId": t.sessionID,
	})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/lasr/run")
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
	resData := startResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return e
	}
	if resData.Code != 0 {
		return errors.New(resData.Desc)
	}
	t.taskID = resData.Data.TaskID
	return nil
}

type taskInfoResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Data struct {
		Progress int `json:"progress"`
	} `json:"data"`
}

func (t *Transcription) GetTaskInfo() (int, error) {
	if t.taskID == "" {
		return 0, errors.New("task uncreated")
	}
	httpClient := t.createhttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"task_id":     t.taskID,
		"x-sessionId": t.sessionID,
	})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/lasr/progress")
	if e != nil {
		return 0, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return 0, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return 0, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return 0, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return 0, errors.New(msg)
	}
	resData := taskInfoResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return 0, e
	}
	if resData.Code != 0 {
		return 0, errors.New(resData.Desc)
	}
	return resData.Data.Progress, nil
}

type transcriptionResponse struct {
	Code int    `json:"code"`
	Desc string `json:"desc"`
	Data struct {
		Result []TranscriptionData `json:"result"`
	} `json:"data"`
}

type TranscriptionData struct {
	Ed      int    `json:"ed"`
	Onebest string `json:"onebest"`
	Bg      int    `json:"bg"`
}

func (t *Transcription) GetResult() ([]TranscriptionData, error) {
	if t.taskID == "" {
		return []TranscriptionData{}, errors.New("task uncreated")
	}
	httpClient := t.createhttpClient()
	httpClient.Header.Set("Content-Type", "application/json")
	httpClient.SetBody(map[string]interface{}{
		"task_id":     t.taskID,
		"x-sessionId": t.sessionID,
	})
	httpRes, e := httpClient.Post("http://api-ai.vivo.com.cn/lasr/result")
	if e != nil {
		return []TranscriptionData{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return []TranscriptionData{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return []TranscriptionData{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return []TranscriptionData{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return []TranscriptionData{}, errors.New(msg)
	}
	resData := transcriptionResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return []TranscriptionData{}, e
	}
	if resData.Code != 0 {
		return []TranscriptionData{}, errors.New(resData.Desc)
	}
	return resData.Data.Result, nil
}
