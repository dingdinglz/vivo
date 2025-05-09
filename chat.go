package vivo

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	CHAT_ROLE_USER      = "user"
	CHAT_ROLE_ASSISTANT = "assistant"
	CHAT_ROLE_SYSTEM    = "system"
	CHAT_ROLE_FUNCTION  = "function"
)

type chatRequest struct {
	Prompt       string        `json:"prompt,omitempty"`
	Messages     []ChatMessage `json:"messages,omitempty"`
	Model        string        `json:"model"`
	SessionID    string        `json:"sessionId"`
	SystemPrompt string        `json:"systemPrompt,omitempty"`
	Extra        *ChatExtra    `json:"extra,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatExtra struct {
	Temperature       float32 `json:"temperature,omitempty"`
	TopP              float32 `json:"top_p,omitempty"`
	TopK              int     `json:"top_k,omitempty"`
	MaxNewTokens      int     `json:"max_new_tokens,omitempty"`
	RepetitionPenalty float32 `json:"repetition_penalty,omitempty"`
}

type chatResponse struct {
	Code int `json:"code"`
	Data struct {
		SessionID        string      `json:"sessionId"`
		RequestID        string      `json:"requestId"`
		Content          string      `json:"content"`
		ReasoningContent interface{} `json:"reasoningContent"`
		Image            interface{} `json:"image"`
		FunctionCall     interface{} `json:"functionCall"`
		ToolCall         interface{} `json:"toolCall"`
		ToolCalls        interface{} `json:"toolCalls"`
		ContentList      interface{} `json:"contentList"`
		SearchInfo       interface{} `json:"searchInfo"`
		Usage            struct {
			PromptTokens     interface{} `json:"promptTokens"`
			CompletionTokens interface{} `json:"completionTokens"`
			TotalTokens      interface{} `json:"totalTokens"`
			Duration         interface{} `json:"duration"`
			ImageCost        interface{} `json:"imageCost"`
			InputImages      interface{} `json:"inputImages"`
			CostLevel        interface{} `json:"costLevel"`
		} `json:"usage"`
		Provider     string      `json:"provider"`
		ClearHistory interface{} `json:"clearHistory"`
		SearchExtra  interface{} `json:"searchExtra"`
		Model        string      `json:"model"`
		FinishReason string      `json:"finishReason"`
		Score        float64     `json:"score"`
		ModelInfo    struct {
			Model        string `json:"model"`
			ModelVersion string `json:"modelVersion"`
		} `json:"modelInfo"`
	} `json:"data"`
	Msg string `json:"msg"`
}

func (app *Vivo) Chat(requestID string, sessionID string, messages []ChatMessage, extra *ChatExtra) (ChatMessage, error) {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", requestID)
	client.Header.Set("Content-Type", "application/json")
	client.SetBody(chatRequest{
		Messages:  messages,
		Model:     "vivo-BlueLM-TB-Pro",
		SessionID: sessionID,
		Extra:     extra,
	})
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/vivogpt/completions")
	if e != nil {
		return ChatMessage{}, e
	}
	defer httpRes.Body.Close()
	body, e := io.ReadAll(httpRes.Body)
	if e != nil {
		return ChatMessage{}, e
	}
	if httpRes.StatusCode() != 200 {
		resMap := make(map[string]interface{})
		e := json.Unmarshal(body, &resMap)
		if e != nil {
			return ChatMessage{}, errors.New(string(body))
		}
		msg, ok := resMap["message"].(string)
		if ok {
			return ChatMessage{}, errors.New(msg)
		}
		msg, _ = resMap["msg"].(string)
		return ChatMessage{}, errors.New(msg)
	}
	resData := chatResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return ChatMessage{}, e
	}
	if resData.Code != 0 {
		return ChatMessage{}, errors.New(resData.Msg)
	}
	return ChatMessage{
		Role:    CHAT_ROLE_ASSISTANT,
		Content: resData.Data.Content,
	}, nil
}

type chatStreamResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Reply   string `json:"reply"`
	Msg     string `json:"msg"`
}

func (app *Vivo) ChatStream(requestID string, sessionID string, messages []ChatMessage, extra *ChatExtra, during func(s string)) error {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", requestID)
	client.Header.Set("Content-Type", "application/json")
	client.SetBody(chatRequest{
		Messages:  messages,
		Model:     "vivo-BlueLM-TB-Pro",
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

func (app *Vivo) EasyChat(sessionID string, message string, systemPrompt ...string) (string, error) {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", GenerateRequestID())
	client.Header.Set("Content-Type", "application/json")
	sysprompt := ""
	if len(systemPrompt) > 0 {
		sysprompt = systemPrompt[0]
	}
	client.SetBody(chatRequest{
		Prompt:       message,
		SystemPrompt: sysprompt,
		Model:        "vivo-BlueLM-TB-Pro",
		SessionID:    sessionID,
	})
	httpRes, e := client.Post("https://api-ai.vivo.com.cn/vivogpt/completions")
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
	resData := chatResponse{}
	e = json.Unmarshal(body, &resData)
	if e != nil {
		return "", e
	}
	if resData.Code != 0 {
		return "", errors.New(resData.Msg)
	}
	return resData.Data.Content, nil
}

func (app *Vivo) EasyChatStream(sessionID string, message string, during func(s string), systemPrompt ...string) error {
	client := app.newHttpClient()
	client.QueryParams.Add("requestId", GenerateRequestID())
	client.Header.Set("Content-Type", "application/json")
	sysprompt := ""
	if len(systemPrompt) > 0 {
		sysprompt = systemPrompt[0]
	}
	client.SetBody(chatRequest{
		Prompt:       message,
		Model:        "vivo-BlueLM-TB-Pro",
		SessionID:    sessionID,
		SystemPrompt: sysprompt,
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
