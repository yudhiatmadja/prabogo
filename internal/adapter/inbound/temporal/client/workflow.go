package client_temporal_inbound_adapter

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"prabogo/internal/model"
)

type ClientWorkflow interface {
	UpsertClientWorkflow(ctx workflow.Context, input model.ClientInput) (string, error)
}

type clientWorkflow struct {
	activity ClientActivity
}

func NewClientWorkflow(
	activity ClientActivity,
) ClientWorkflow {
	return &clientWorkflow{
		activity: activity,
	}
}

func (g *clientWorkflow) UpsertClientWorkflow(ctx workflow.Context, input model.ClientInput) (string, error) {
	logger := workflow.GetLogger(ctx)
	workflowInfo := workflow.GetInfo(ctx)

	logger.Info("workflow started", "workflowID", workflowInfo.WorkflowExecution.ID)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var results []model.Client
	err := workflow.ExecuteActivity(
		ctx,
		g.activity.Upsert,
		[]model.ClientInput{input},
	).Get(ctx, &results)
	if err != nil {
		logger.Error("upsert client activity failed", "error", err)
		return "failed to upsert client", err
	}

	var bearerKey string
	if len(results) > 0 {
		bearerKey = results[0].BearerKey
	}

	successMessage := "bearer key: " + bearerKey
	logger.Info(successMessage, "workflowID", workflowInfo.WorkflowExecution.ID)

	return successMessage, nil
}
