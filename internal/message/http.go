package message

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client interface {
	SendMessage(ctx context.Context, number string, message string) error
}

const (
	wppUrl = "https://graph.facebook.com/v13.0/%s/messages"
)

var (
	errGenericClient = errors.New("error sending message to whatsapp")
)

type client struct {
	PhoneID    string
	HttpClient *http.Client
	Token      string
}

func (c client) SendMessage(ctx context.Context, number string, message string) error {
	log.Printf("[INFO] starting send message client | number: %s | message: %s", number, message)

	body := buildMessageBody(number, message)
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf(wppUrl, c.PhoneID), bytes.NewBuffer(jsonBody))
	c.addHeaders(req)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] Error message client - http | error: %s", err.Error())
		return err
	}
	defer resp.Body.Close()

	responseJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Error message client - read body | error: %s", err.Error())
		return err
	}

	var response WppMessageResponse
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		log.Printf("[ERROR] Error track client - unmarshal | error: %s", err.Error())
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Error track client - status code | response_code: %d | response_body: %s", resp.StatusCode, string(responseJson))
		return errGenericClient
	}

	log.Printf("[INFO] success send message client | response_code: %d | response_body: %s", resp.StatusCode, string(responseJson))
	return nil
}

func (c client) addHeaders(req *http.Request) {
	req.Header.Set("Authorization", c.Token)
	req.Header.Set("Content-Type", "application/json")
}

func buildMessageBody(number string, message string) WppMessageRequest {
	return WppMessageRequest{
		MessagingProduct: "whatsapp",
		To:               number,
		Text:             Text{Body: message},
	}
}

func NewClient() Client {
	return &client{
		PhoneID:    "PHONE_ID",
		HttpClient: &http.Client{},
		Token:      "TOKEN",
	}
}
