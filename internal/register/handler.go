package register

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"wpp-cloud/internal/domain"
)

type Handler interface {
	RegisterWebhook(context.Context, domain.Request) (domain.Response, error)
}

type handler struct {
	Token string `json:"token"`
}

func (h *handler) RegisterWebhook(ctx context.Context, request domain.Request) (domain.Response, error) {
	mode := request.QueryParams["hub.mode"]
	token := request.QueryParams["hub.verify_token"]
	challenge := request.QueryParams["hub.challenge"]

	challengeInt, err := strconv.ParseInt(challenge, 10, 64)
	if err != nil {
		return domain.Response{}, err
	}

	if len(mode) != 0 && mode == "subscribe" && token == h.Token {
		log.Printf("[INFO] success registering webhook - challenge: %d", challengeInt)
		return domain.Response{Status: http.StatusOK, Body: challengeInt}, nil
	}

	return domain.Response{}, domain.NewUnauthorizedError("invalid request")
}

func NewHandler(token string) Handler {
	return &handler{Token: token}
}
