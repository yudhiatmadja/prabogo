package mcp_utils

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/palantir/stacktrace"

	"prabogo/internal/model"
	"prabogo/utils/log"
)

func RunHTTPServer(ctx context.Context, server *mcp.Server) {
	addr := os.Getenv("MCP_HTTP_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{JSONResponse: true, Stateless: true})

	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.WithContext(ctx).Infof("mcp server started on http://localhost%s/mcp", addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithContext(ctx).Fatalf("failed to run mcp http server: %+v", err)
	}
}

func RunStdioServer(ctx context.Context, server *mcp.Server) {
	log.WithContext(ctx).Info("mcp stdio server started")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.WithContext(ctx).Fatalf("failed to run mcp stdio server: %+v", err)
	}
}

func AddTool[T any](data []T, name string, err error) (*mcp.CallToolResult, model.Response, error) {
	if err != nil {
		return nil, model.Response{}, fmt.Errorf("mcp get_%s failed: %s", name, stacktrace.RootCause(err).Error())
	}

	return SummaryToolResult(name, len(data)), model.Response{Success: true, Data: data}, nil
}

func SummaryToolResult(objectType string, objectLen int) *mcp.CallToolResult {
	summary := fmt.Sprintf("No %ss found", objectType)
	if objectLen > 0 {
		summary = fmt.Sprintf("Found %d %s", objectLen, objectType)
		if objectLen != 1 {
			summary += "s"
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}
}
