package tools

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// jsonResult marshals v to a JSON TextContent result.
func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errorResult(fmt.Errorf("marshal result: %w", err))
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(b)},
		},
	}, nil, nil
}

// textResult wraps a plain string in a TextContent result.
func textResult(s string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: s},
		},
	}, nil, nil
}

// errorResult returns a tool-execution error as an isError:true result.
// Per MCP spec §6, API/business failures must be reported this way rather
// than as protocol-level Go errors so the client can distinguish them.
func errorResult(err error) (*mcp.CallToolResult, any, error) {
	msg := "unknown error"
	if err != nil {
		msg = err.Error()
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
		IsError: true,
	}, nil, nil
}
