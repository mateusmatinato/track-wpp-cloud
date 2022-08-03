package webhook

import (
	"context"
	"log"
	"wpp-cloud/internal/commands"
)

const (
	TypeMessages = "message-received"
	TypeStatus   = "message-status"
	MessageText  = "text"
	UnknownType  = ""
)

type Service interface {
	ProcessWebhook(ctx context.Context, request WebhookRequestDTO) error
}

type service struct {
	commandService commands.Service
}

func (s service) ProcessWebhook(ctx context.Context, request WebhookRequestDTO) error {
	log.Printf("[INFO] starting processing webhook, request: %v", request)

	webhookType := getWebhookType(request)
	if webhookType == UnknownType {
		log.Printf("[WARN] webhook have a unknown type, %s", webhookType)
		return nil
	}

	switch webhookType {
	case TypeMessages:
		mainObj := request.Entry[0].Changes[0].Value
		if mainObj.Contacts == nil {
			log.Printf("[ERROR] webhook doesnt have a contact")
		}
		contacts := *mainObj.Contacts
		return s.ProcessMessageWebhook(ctx, *mainObj.Messages, contacts[0])
	case TypeStatus:
		log.Printf("[WARN] currently, message-status webhook are only logged")
	}

	return nil
}

func NewService(commandService commands.Service) Service {
	return &service{commandService: commandService}
}

func getWebhookType(dto WebhookRequestDTO) string {
	mainObj := dto.Entry[0].Changes[0].Value
	if mainObj.Messages != nil {
		return TypeMessages
	} else if mainObj.Statuses != nil {
		return TypeStatus
	} else {
		return UnknownType
	}
}

func (s service) ProcessMessageWebhook(ctx context.Context, messages []MessageWebhook, contact ContactWebhook) error {
	message := messages[0]
	log.Printf("[INFO] starting process message webhook | type: %s | contact: %s", message.Type, contact)

	switch message.Type {
	case MessageText:
		log.Printf("[INFO] received text message | message: %s", *message.Text)

		commandRequest := commands.ProcessRequest{
			Message: message.Text.Body,
			Contact: commands.Contact{
				Number: contact.WaID,
				Name:   contact.Profile.Name,
			},
		}
		return s.commandService.ProcessMessage(ctx, commandRequest)
	default:
		log.Printf("[ERROR] received message from a format that is currently not suported | type: %s | message: %s | contact: %s", message.Type, message, contact)
	}

	return nil
}
