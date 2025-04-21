package vivo

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type websocketError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (app *Vivo) newWebsocketClient(u url.URL) (*websocket.Conn, error) {
	header := http.Header{}
	header.Add("X-AI-GATEWAY-APP-ID", app.appID)
	header.Add("X-AI-GATEWAY-TIMESTAMP", int64toString(time.Now().Unix()))
	header.Add("X-AI-GATEWAY-NONCE", generateRandomString(8))
	header.Add("X-AI-GATEWAY-SIGNED-HEADERS", "x-ai-gateway-app-id;x-ai-gateway-timestamp;x-ai-gateway-nonce")
	signing_string := "GET" + "\n" + u.Path + "\n" + u.Query().Encode() + "\n" + app.appID + "\n" + header.Get("X-AI-GATEWAY-TIMESTAMP") + "\n" + "x-ai-gateway-app-id" + ":" + app.appID + "\n" + "x-ai-gateway-timestamp" + ":" + header.Get("X-AI-GATEWAY-TIMESTAMP") + "\n" + "x-ai-gateway-nonce" + ":" + header.Get("X-AI-GATEWAY-NONCE")
	header.Set("X-AI-GATEWAY-SIGNATURE", base64encode(hMACSHA256HEX(signing_string, app.appKey)))
	c, http, e := websocket.DefaultDialer.Dial(u.String(), header)
	defer http.Body.Close()
	body, _ := io.ReadAll(http.Body)
	data := websocketError{}
	json.Unmarshal(body, &data)
	if data.ErrorCode != 0 {
		return nil, errors.New(data.ErrorMsg)
	}
	return c, e
}
