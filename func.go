package vivo

import (
	"encoding/json"
	"errors"
	"strings"
)

type ChatTool struct {
	FuncName    string
	Description string
	Parameters  []ChatToolParameter
	Func        func(map[string]interface{}) (string, error)
}

type rawChatTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  struct {
		Type       string                       `json:"type"`
		Properties map[string]ChatToolParameter `json:"properties"`
		Required   []string                     `json:"required"`
	} `json:"parameters"`
}

type ChatToolParameter struct {
	Name        string   `json:"-"`
	Type        string   `json:"type"`
	Enum        []string `json:"enum,omitempty"`
	Description string   `json:"description"`
	Required    bool     `json:"-"`
}

type chatToolCall struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

func makeRawTools(tools []ChatTool) string {
	rawResult := make([]rawChatTool, 0)
	for _, item := range tools {
		newRawTool := rawChatTool{}
		newRawTool.Name = item.FuncName
		newRawTool.Description = item.Description
		newRawTool.Parameters.Type = "object"
		newRawTool.Parameters.Properties = make(map[string]ChatToolParameter)
		requiredParameters := make([]string, 0)
		for _, item2 := range item.Parameters {
			newRawTool.Parameters.Properties[item2.Name] = item2
			if item2.Required {
				requiredParameters = append(requiredParameters, item2.Name)
			}
		}
		newRawTool.Parameters.Required = requiredParameters
		rawResult = append(rawResult, newRawTool)
	}
	res, _ := json.MarshalIndent(rawResult, "", "    ")
	return string(res)
}

func (app *Vivo) ChatWithTools(sessionID string, messages []ChatMessage, extra *ChatExtra, tools []ChatTool) (string, error) {
	if len(tools) == 0 {
		return "", errors.New("至少携带一个工具")
	}
	if len(messages) == 0 {
		return "", errors.New("消息不能为空")
	}
	if messages[len(messages)-1].Role == CHAT_ROLE_USER {
		messages = append([]ChatMessage{
			{
				Role: CHAT_ROLE_SYSTEM,
				Content: `你是一个AI助手，尽你所能回答用户的问题。

你可以使用的工具如下:
<APIS>
` + makeRawTools(tools) + "\n" + `</APIS>

如果用户的问题需要调用工具，输出格式为：
<APIs>
[{"name": "函数名","parameters": {"参数名": "参数"}}]
</APIs>
否则直接回复用户。`,
			},
		}, messages...)
	}
	resMessage, e := app.Chat(GenerateRequestID(), sessionID, messages, extra)
	if e != nil {
		return "", e
	}
	if strings.HasPrefix(resMessage.Content, "<APIs>") && strings.HasSuffix(resMessage.Content, "</APIs>") {
		messages = append(messages, ChatMessage{
			Role:    CHAT_ROLE_ASSISTANT,
			Content: resMessage.Content,
		})
		ToolCallsMessage := resMessage.Content[6 : len(resMessage.Content)-7]
		var toolCalls []chatToolCall
		e := json.Unmarshal([]byte(ToolCallsMessage), &toolCalls)
		if e != nil {
			return "", errors.New("tool call error")
		}
		var resDatas []string
		for _, item := range toolCalls {
			flag := false
			for _, item2 := range tools {
				if item2.FuncName == item.Name {
					toolCallRes, e := item2.Func(item.Parameters)
					if e != nil {
						return "", errors.New("tool " + item.Name + ":" + e.Error())
					}
					resDatas = append(resDatas, toolCallRes)
					flag = true
					break
				}
			}
			if !flag {
				resDatas = append(resDatas, "")
			}
		}
		sendToolCallRes, _ := json.Marshal(resDatas)
		messages = append(messages, ChatMessage{
			Role:    CHAT_ROLE_FUNCTION,
			Content: string(sendToolCallRes),
		})
		res, e := app.ChatWithTools(sessionID, messages, extra, tools)
		return res, e
	}
	return resMessage.Content, nil
}
