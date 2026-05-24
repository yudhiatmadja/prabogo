package inbound_port

type McpPort interface {
	Client() ClientMcpPort
}
