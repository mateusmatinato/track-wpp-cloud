package recurring

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

type Handler interface {
	ReceiveWebhook(context.Context, events.CloudWatchEvent) error
}

type handler struct {
	Service Service
}

func (h handler) ReceiveWebhook(ctx context.Context, req events.CloudWatchEvent) error {
	fmt.Printf("[INFO] Starting recurring event handler | event_id: %s", req.ID)
	return h.Service.ProcessCodes(ctx)
}

func NewHandler(service Service) Handler {
	return &handler{Service: service}
}
