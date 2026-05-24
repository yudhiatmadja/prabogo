package rabbitmq_inbound_adapter

import (
	"context"
	"encoding/json"

	"github.com/palantir/stacktrace"

	"prabogo/internal/domain"
	"prabogo/internal/model"
	inbound_port "prabogo/internal/port/inbound"
	"prabogo/utils/activity"
	"prabogo/utils/log"
)

type clientAdapter struct {
	domain domain.Domain
}

func NewClientAdapter(
	domain domain.Domain,
) inbound_port.ClientMessagePort {
	return &clientAdapter{
		domain: domain,
	}
}

func (h *clientAdapter) Upsert(a any) bool {
	msg := a.([]byte)
	ctx := activity.NewContext("message_client_upsert")
	var payload []model.ClientInput
	err := json.Unmarshal(msg, &payload)
	if err != nil {
		log.WithContext(ctx).Errorf("message client upsert failed: %s", err.Error())
		return true
	}
	ctx = context.WithValue(ctx, activity.Payload, payload)

	results, err := h.domain.Client().Upsert(ctx, payload)
	if err != nil {
		log.WithContext(ctx).Errorf("message client upsert failed: %s", stacktrace.RootCause(err).Error())
	}
	ctx = context.WithValue(ctx, activity.Result, results)

	log.WithContext(ctx).Info("message client upsert success")
	return true
}
