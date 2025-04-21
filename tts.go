package vivo

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	TTS_MODE_SHORT = "short_audio_synthesis_jovi"
	TTS_MODE_LONG  = "long_audio_synthesis_screen"
	TTS_MODE_HUMAN = "tts_humanoid_lam"
)

type ttsRequest struct {
	Aue      int    `json:"aue"`
	Auf      string `json:"auf"`
	Vcn      string `json:"vcn"`
	Speed    int    `json:"speed,omitempty"`
	Volume   int    `json:"volume,omitempty"`
	Text     string `json:"text"`
	Encoding string `json:"encoding"`
	ReqID    int    `json:"reqId,omitempty"`
}

type ttsResponse struct {
	Data struct {
		Audio  string `json:"audio"`
		Status int    `json:"status"`
	} `json:"data"`
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (app *Vivo) TTS(mode string, vcn string, text string) ([]byte, error) {
	u := url.URL{Scheme: "wss", Host: "api-ai.vivo.com.cn", Path: "/tts"}
	query := u.Query()
	query.Set("engineid", mode)
	query.Set("system_time", int64toString(time.Now().Unix()))
	query.Set("user_id", "userX")
	query.Set("model", "modelX")
	query.Set("product", "productX")
	query.Set("package", "packageX")
	query.Set("client_version", "0")
	query.Set("system_version", "0")
	query.Set("sdk_version", "0")
	query.Set("android_version", "9")

	u.RawQuery = query.Encode()
	conn, e := app.newWebsocketClient(u)
	if e != nil {
		return []byte{}, e
	}
	wg := sync.WaitGroup{}
	defer conn.Close()
	wg.Add(1)
	var innerError error
	var pcmData []byte
	go func() {
		defer wg.Done()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				innerError = err
				return
			}
			resData := ttsResponse{}
			json.Unmarshal(message, &resData)
			if resData.ErrorCode != 0 {
				innerError = errors.New(resData.ErrorMsg)
				return
			}
			appendData, err := base64.StdEncoding.DecodeString(resData.Data.Audio)
			if err != nil {
				innerError = err
				return
			}
			pcmData = append(pcmData, appendData...)
			if resData.Data.Status == 2 {
				return
			}
		}
	}()
	sendData, _ := json.Marshal(ttsRequest{
		Auf:      "audio/L16;rate=24000",
		Vcn:      vcn,
		Text:     base64encode([]byte(text)),
		Encoding: "utf8",
		ReqID:    513722013,
	})
	conn.WriteMessage(websocket.TextMessage, sendData)
	wg.Wait()
	if innerError != nil {
		return []byte{}, innerError
	}
	return pcmData, nil
}
