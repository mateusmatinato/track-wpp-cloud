package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"wpp-cloud/internal/message"
	"wpp-cloud/internal/recurring"
	"wpp-cloud/internal/repository"
	"wpp-cloud/internal/tracking"
)

func main() {

	trackerClient := tracking.NewClient()
	trackerService := tracking.NewService(trackerClient)

	messageClient := message.NewClient()
	messageService := message.NewService(messageClient)

	repo := repository.NewRepository()

	recurringService := recurring.NewService(trackerService, messageService, repo)
	recurringHandler := recurring.NewHandler(recurringService)

	api := recurringAPI{
		RecurringHandler: recurringHandler,
	}

	lambda.Start(api.router)
}

type recurringAPI struct {
	RecurringHandler recurring.Handler
}

func (h *recurringAPI) router(ctx context.Context, req events.CloudWatchEvent) error {
	log.Printf("[INFO] starting cloudwatch event | event_id: %s | event_detail: %s", req.ID, req.Detail)
	return h.RecurringHandler.ReceiveWebhook(ctx, req)
}
