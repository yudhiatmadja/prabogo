package mcp_inbound_adapter

import (
	"context"

	"github.com/palantir/stacktrace"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
	"prabogo/utils/activity"
)

type clientAdapter struct {
	domain domain.Domain
}

func NewClientAdapter(
	domain domain.Domain,
) inbound_port.ClientMcpPort {
	return &clientAdapter{
		domain: domain,
	}
}

func (a *clientAdapter) Find(payload any) ([]model.Client, error) {
	filter, ok := payload.(model.ClientFilter)
	if !ok {
		return nil, stacktrace.NewError("invalid payload type: expected model.ClientFilter")
	}

	ctx := activity.NewContext("mcp_client_find_by_filter")
	ctx = context.WithValue(ctx, activity.Payload, filter)

	results, err := a.domain.Client().FindByFilter(ctx, filter)
	if err != nil {
		return nil, stacktrace.Propagate(err, "find client by filter error")
	}

	return results, nil
}
