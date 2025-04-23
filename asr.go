package vivo

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/gorilla/websocket"
)

type asrShortVoiceRecognitionResponse struct {
	Data struct {
		IsLast bool   `json:"is_last"`
		Text   string `json:"text"`
	} `json:"data"`
	Code int    `json:"code"`
	Desc string `json:"desc"`
}

func (app *Vivo) AsrShortVoiceRecognition(file string) (string, error) {
	wavFile, e := os.Open(file)
	if e != nil {
		return "", e
	}
	defer wavFile.Close()

	decoder := wav.NewDecoder(wavFile)
	if !decoder.IsValidFile() {
		return "", errors.New("not a valid wave file")
	}

	decoder.ReadInfo()

	e = decoder.FwdToPCM()
	if e != nil {
		return "", e
	}

	u := url.URL{Scheme: "wss", Host: "api-ai.vivo.com.cn", Path: "/asr/v2"}
	query := u.Query()
	query.Set("engineid", "shortasrinput")
	query.Set("system_time", int64toString(time.Now().Unix()))
	query.Set("user_id", "userX")
	query.Set("model", "modelX")
	query.Set("package", "packageX")
	query.Set("client_version", "0")
	query.Set("system_version", "0")
	query.Set("sdk_version", "0")
	query.Set("android_version", "9")
	query.Set("net_type", "0")

	u.RawQuery = query.Encode()
	conn, e := app.newWebsocketClient(u)
	if e != nil {
		return "", e
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var innerError error
	var resString string
	go func() {
		defer wg.Done()
		for {
			_, message, e := conn.ReadMessage()
			if e != nil {
				return
			}
			data := asrShortVoiceRecognitionResponse{}
			json.Unmarshal(message, &data)
			if data.Code != 0 {
				conn.WriteMessage(websocket.BinaryMessage, []byte("--close--"))
				innerError = errors.New(data.Desc)
				return
			}
			resString = data.Data.Text
			if data.Data.IsLast {
				conn.WriteMessage(websocket.BinaryMessage, []byte("--close--"))
				return
			}
		}
	}()

	conn.WriteJSON(map[string]interface{}{
		"type":       "started",
		"request_id": GenerateRequestID(),
		"asr_info": map[string]interface{}{
			"end_vad_time":    100,
			"audio_type":      "pcm",
			"punctuation":     1,
			"chinese2digital": 0,
		},
	})

	buf := &audio.IntBuffer{
		Data: make([]int, 4096),
	}

	for {
		n, e := decoder.PCMBuffer(buf)
		if e != nil || n == 0 {
			break
		}
		pcmBytes, e := pcmIntToBytes(buf.Data[:n], int(decoder.BitDepth))
		if e != nil || n == 0 {
			break
		}
		conn.WriteMessage(websocket.BinaryMessage, pcmBytes)
	}
	conn.WriteMessage(websocket.BinaryMessage, []byte("--end--"))

	wg.Wait()

	if innerError != nil {
		return "", nil
	}

	return resString, nil
}

type asrLongVoiceRecognitionResponse struct {
	Code int `json:"code"`
	Data struct {
		Onebest string `json:"onebest"`
	} `json:"data"`
	Desc string `json:"desc"`
}

func (app *Vivo) AsrLongVoiceRecognition(file string) (string, error) {
	wavFile, e := os.Open(file)
	if e != nil {
		return "", e
	}
	defer wavFile.Close()

	decoder := wav.NewDecoder(wavFile)
	if !decoder.IsValidFile() {
		return "", errors.New("not a valid wave file")
	}

	decoder.ReadInfo()

	e = decoder.FwdToPCM()
	if e != nil {
		return "", e
	}

	u := url.URL{Scheme: "wss", Host: "api-ai.vivo.com.cn", Path: "/asr/v2"}
	query := u.Query()
	query.Set("engineid", "longasrlisten")
	query.Set("system_time", int64toString(time.Now().Unix()))
	query.Set("user_id", "userX")
	query.Set("model", "modelX")
	query.Set("package", "packageX")
	query.Set("product", "productX")
	query.Set("client_version", "0")
	query.Set("system_version", "0")
	query.Set("sdk_version", "0")
	query.Set("android_version", "9")
	query.Set("net_type", "0")

	u.RawQuery = query.Encode()
	conn, e := app.newWebsocketClient(u)
	if e != nil {
		return "", e
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var innerError error
	var resString string
	go func() {
		defer wg.Done()
		for {
			_, message, e := conn.ReadMessage()
			if e != nil {
				return
			}
			data := asrLongVoiceRecognitionResponse{}
			json.Unmarshal(message, &data)
			if data.Code != 0 && data.Code != 8 && data.Code != 9 {
				innerError = errors.New(data.Desc)
				conn.WriteMessage(websocket.BinaryMessage, []byte("--close--"))
				return
			}
			if data.Code == 9 {
				resString = data.Data.Onebest
				conn.WriteMessage(websocket.BinaryMessage, []byte("--close--"))
				return
			}
		}
	}()

	conn.WriteJSON(map[string]interface{}{
		"type":       "started",
		"request_id": GenerateRequestID(),
		"asr_info": map[string]interface{}{
			"end_vad_time": 100,
			"audio_type":   "pcm",
			"punctuation":  1,
		},
	})

	buf := &audio.IntBuffer{
		Data: make([]int, 4096),
	}

	for {
		n, e := decoder.PCMBuffer(buf)
		if e != nil || n == 0 {
			break
		}
		pcmBytes, e := pcmIntToBytes(buf.Data[:n], int(decoder.BitDepth))
		if e != nil || n == 0 {
			break
		}
		conn.WriteMessage(websocket.BinaryMessage, pcmBytes)
	}
	conn.WriteMessage(websocket.BinaryMessage, []byte("--end--"))

	wg.Wait()

	if innerError != nil {
		return "", nil
	}

	return resString, nil
}
