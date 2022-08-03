package message

import (
	"context"
	"log"
)

type Service interface {
	SendMessage(ctx context.Context, number string, message string) error
}

type service struct {
	Client Client
}

func (s service) SendMessage(ctx context.Context, number string, message string) error {
	log.Printf("[INFO] starting send message service | number: %s | message: %s", number, message)

	err := s.Client.SendMessage(ctx, number, message)
	return err
}

func NewService(client Client) Service {
	return &service{Client: client}
}
