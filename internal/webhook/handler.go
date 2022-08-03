package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"wpp-cloud/internal/domain"
)

type Handler interface {
	ReceiveWebhook(context.Context, domain.Request) (domain.Response, error)
}

type handler struct {
	FbSecret string
	Service
}

func (h *handler) ReceiveWebhook(ctx context.Context, request domain.Request) (domain.Response, error) {
	bodySignature := request.Headers["X-Hub-Signature-256"]

	err := h.validateBodySignature(bodySignature, request.Body)
	if err != nil {
		return domain.Response{}, err
	}

	var dto WebhookRequestDTO
	err = json.Unmarshal([]byte(request.Body), &dto)
	if err != nil {
		log.Printf("[ERROR] error unmarshalling body, error: %s", err.Error())
		return domain.Response{}, err
	}

	err = h.ProcessWebhook(ctx, dto)
	if err != nil {
		return domain.Response{}, err
	}

	return domain.Response{
		Status: 200,
		Body:   "success processing webhook",
	}, nil
}

func (h *handler) validateBodySignature(signature string, body string) error {
	hash := hmac.New(sha256.New, []byte(h.FbSecret))
	hash.Write([]byte(body))
	sha := "sha256=" + hex.EncodeToString(hash.Sum(nil))
	if sha != signature {
		log.Printf("[ERROR] error validating body signature")
		return domain.NewInternalServerError(fmt.Sprintf("error validating body signature, signatureReceived:%s, signatureCalculated:%s", signature, sha))
	}
	return nil
}

func NewHandler(FbSecret string, service Service) Handler {
	return &handler{
		FbSecret: FbSecret,
		Service:  service,
	}
}
