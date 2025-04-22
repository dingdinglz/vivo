package vivo

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	CHAT_MESSAGE_TEXT  = "text"
	CHAT_MESSAGE_IMAGE = "image"

	VISION_CHAT_MODEL_BLUELM_PRD = "BlueLM-Vision-prd"
	VISION_CHAT_MODEL_BLUELM_V2  = "vivo-BlueLM-V-2.0"
)

type visionChatRequest struct {
	Prompt       string              `json:"prompt,omitempty"`
	Messages     []VisionChatMessage `json:"messages,omitempty"`
	Model        string              `json:"model"`
	SessionID    string              `json:"sessionId"`
	SystemPrompt string              `json:"systemPrompt,omitempty"`
	Extra        *ChatExtra          `json:"extra,omitempty"`
}

type VisionChatMessage struct {
	Role        string `json:"role"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

func (app *Vivo) VisionChat(requestID string, sessionID string, model string, messages []VisionChatMessage, extra *ChatExtra) (VisionChatMessage, error) {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", requestID)
	client.Header.Set("Content-Type", "application/json")
	client.SetBody(visionChatRequest{
		Messages:  messages,
		Model:     model,
		SessionID: sessionID,
		Extra:     extra,
	})
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/vivogpt/completions")
	if e != nil {
		return VisionChatMessage{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return VisionChatMessage{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return VisionChatMessage{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return VisionChatMessage{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return VisionChatMessage{}, errors.New(msg)
	}
	resData := chatResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return VisionChatMessage{}, e
	}
	if resData.Code != 0 {
		return VisionChatMessage{}, errors.New(resData.Msg)
	}
	return VisionChatMessage{
		Role:        CHAT_ROLE_ASSISTANT,
		ContentType: CHAT_MESSAGE_TEXT,
		Content:     resData.Data.Content,
	}, nil
}

func (app *Vivo) VisionChatStream(requestID string, sessionID string, model string, messages []VisionChatMessage, extra *ChatExtra, during func(s string)) error {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", requestID)
	client.Header.Set("Content-Type", "application/json")
	client.SetBody(visionChatRequest{
		Messages:  messages,
		Model:     model,
		SessionID: sessionID,
		Extra:     extra,
	})
	client.SetDoNotParseResponse(true)
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/vivogpt/completions/stream")
	if e != nil {
		return e
	}
	defer httpRes.Body.Close()
	scanner := bufio.NewScanner(httpRes.Body)
	errorEvent := false
	antispamEvent := false
	for scanner.Scan() {
		res := scanner.Text()
		if res == "" {
			continue
		}
		if strings.HasPrefix(res, "data:") {
			// 正常对话
			data := res[5:]
			resData := chatStreamResponse{}
			json.Unmarshal([]byte(data), &resData)
			if errorEvent {
				return errors.New(resData.Msg)
			}
			if antispamEvent {
				return errors.New(resData.Reply)
			}
			during(resData.Message)
		}
		if strings.HasPrefix(res, "event:") {
			// 事件
			event := res[6:]
			switch event {
			case "error":
				errorEvent = true
			case "antispam":
				antispamEvent = true
			case "close":
				return nil
			default:
			}
		}
	}
	return nil
}
