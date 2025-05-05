package vivo

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func McpToTools(command string, env []string, args ...string) ([]ChatTool, error) {
	c, e := client.NewStdioMCPClient(command, env, args...)
	if e != nil {
		return []ChatTool{}, e
	}
	ctx := context.Background()
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "vivo-aigc-client",
		Version: "1.0.0",
	}
	_, e = c.Initialize(ctx, initRequest)
	if e != nil {
		return []ChatTool{}, e
	}
	toolsRequest := mcp.ListToolsRequest{}
	tools, e := c.ListTools(ctx, toolsRequest)
	if e != nil {
		return []ChatTool{}, e
	}
	resChatTools := make([]ChatTool, 0)
	for _, item := range tools.Tools {
		newChatTool := ChatTool{}
		newChatTool.FuncName = item.Name
		newChatTool.Description = item.Description
		var newParameters []ChatToolParameter

		for name, prop := range item.InputSchema.Properties {
			propMap, ok := prop.(map[string]interface{})
			if !ok {
				continue
			}
			newParameter := ChatToolParameter{}
			newParameter.Name = name
			newParameter.Type, _ = propMap["type"].(string)
			newParameter.Description, _ = propMap["description"].(string)
			newParameters = append(newParameters, newParameter)
		}
		newChatTool.Parameters = newParameters
		newChatTool.Func = func(m map[string]interface{}) (string, error) {
			listTmpRequest := mcp.CallToolRequest{}
			listTmpRequest.Params.Name = item.Name
			listTmpRequest.Params.Arguments = m
			result, e := c.CallTool(ctx, listTmpRequest)
			resData := ""
			for _, content := range result.Content {
				if resData != "" {
					resData += "\n"
				}
				if textContent, ok := content.(mcp.TextContent); ok {
					resData += textContent.Text
				} else {
					jsonBytes, _ := json.MarshalIndent(content, "", "  ")
					resData += string(jsonBytes)
				}
			}
			return resData, e
		}
		resChatTools = append(resChatTools, newChatTool)
	}
	return resChatTools, nil
}
