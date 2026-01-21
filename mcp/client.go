package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// MCPClient представляет клиент для взаимодействия с MCP Server
type MCPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewMCPClient создает новый MCP клиент
func NewMCPClient(baseURL string) *MCPClient {
	return &MCPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ToolCallRequest представляет запрос на вызов tool через MCP
type ToolCallRequest struct {
	ToolName   string                 `json:"toolName"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ToolCallResponse представляет ответ от MCP Server
type ToolCallResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

// MCPRequest представляет общий запрос к MCP Server
type MCPRequest struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

// MCPResponse представляет общий ответ от MCP Server
type MCPResponse struct {
	Result interface{} `json:"result"`
	Error  *MCPError   `json:"error,omitempty"`
}

// MCPError представляет ошибку от MCP Server
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CallTools вызывает один или несколько tools на MCP Server
func (c *MCPClient) CallTools(ctx context.Context, toolNames []string, parameters map[string]interface{}) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for _, toolName := range toolNames {
		result, err := c.callTool(ctx, toolName, parameters)
		if err != nil {
			results[toolName] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}
		results[toolName] = result
	}

	return results, nil
}

// callTool вызывает один tool на MCP Server
func (c *MCPClient) callTool(ctx context.Context, toolName string, parameters map[string]interface{}) (interface{}, error) {
	// Формируем запрос в формате MCP
	request := MCPRequest{
		Method: "tools/call",
		Params: map[string]interface{}{
			"name":       toolName,
			"arguments":  parameters,
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP Server returned status %d: %s", resp.StatusCode, string(body))
	}

	var mcpResp MCPResponse
	if err := json.Unmarshal(body, &mcpResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if mcpResp.Error != nil {
		return nil, fmt.Errorf("MCP Server error: %s (code: %d)", mcpResp.Error.Message, mcpResp.Error.Code)
	}

	return mcpResp.Result, nil
}
