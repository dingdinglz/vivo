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

func (app *Vivo) DialectRecognition(file string) (string, error) {
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

	u := url.URL{Scheme: "ws", Host: "api-ai.vivo.com.cn", Path: "/asr/v2"}
	query := u.Query()
	query.Set("engineid", "shortasrinput")
	query.Set("system_time", int64toString(time.Now().Unix()))
	query.Set("user_id", "2addc42b7ae689dfdf1c63e220df52a2")
	query.Set("package", "unknown")
	query.Set("client_version", "unknown")
	query.Set("sdk_version", "unknown")
	query.Set("android_version", "unknown")
	query.Set("net_type", "1")
	query.Set("product", "x")
	query.Set("user_info", "0")

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
		"request_id": "req_id",
		"asr_info": map[string]interface{}{
			"front_vad_time":  6000,
			"end_vad_time":    2000,
			"audio_type":      "pcm",
			"chinese2digital": 1,
			"punctuation":     2,
			"lang":            "dialect",
		},
		"business_info": "{\"scenes_pkg\":\"com.tencent.qqlive\", \"editor_type\":\"3\", \"pro_id\":\"2addc42b7ae689dfdf1c63e220df52a2-2020\"}",
	})

	buf := &audio.IntBuffer{
		Data: make([]int, 640),
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
