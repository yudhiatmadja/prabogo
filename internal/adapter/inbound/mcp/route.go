package mcp_inbound_adapter

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
	mcp_utils "prabogo/utils/mcp"
)

func InitRoute(
	ctx context.Context,
	server *mcp.Server,
	port inbound_port.McpPort,
) {
	// Add get clients tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_clients",
		Description: "Get clients from prabogo database by filter",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input model.ClientFilter) (*mcp.CallToolResult, model.Response, error) {
		clients, err := port.Client().Find(input)
		return mcp_utils.AddTool(clients, "client", err)
	})
}
