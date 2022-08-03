package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"wpp-cloud/internal/commands"
	"wpp-cloud/internal/domain"
	"wpp-cloud/internal/message"
	"wpp-cloud/internal/register"
	"wpp-cloud/internal/repository"
	"wpp-cloud/internal/tracking"
	"wpp-cloud/internal/webhook"
)

func main() {

	registerHandler := register.NewHandler("TOKEN")

	trackerClient := tracking.NewClient()
	trackerService := tracking.NewService(trackerClient)

	messageClient := message.NewClient()
	messageService := message.NewService(messageClient)

	repo := repository.NewRepository()

	commandService := commands.NewService(trackerService, messageService, repo)

	webhookService := webhook.NewService(commandService)
	webhookHandler := webhook.NewHandler("SECRET", webhookService)

	api := wppAPI{
		RegisterHandler: registerHandler,
		WebhookHandler:  webhookHandler,
	}

	lambda.Start(api.router)
}

type wppAPI struct {
	RegisterHandler register.Handler
	WebhookHandler  webhook.Handler
}

func (h *wppAPI) router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("[INFO] starting request body: %s, queryParams: %s, headers: %s", req.Body, req.QueryStringParameters, req.Headers)
	request := domain.BuildRequest(req)

	switch req.Path {
	case "/default/webhook":
		if req.HTTPMethod == "GET" {
			return domain.BuildAPIGWResponse(h.RegisterHandler.RegisterWebhook(ctx, *request))
		}

		if req.HTTPMethod == "POST" {
			return domain.BuildAPIGWResponse(h.WebhookHandler.ReceiveWebhook(ctx, *request))
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       http.StatusText(http.StatusMethodNotAllowed),
	}, nil
}
