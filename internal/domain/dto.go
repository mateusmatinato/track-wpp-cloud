package domain

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

type Request struct {
	QueryParams map[string]string `json:"query_params"`
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
}

type Response struct {
	Status int `json:"status"`
	Body   interface{}
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func BuildRequest(req events.APIGatewayProxyRequest) *Request {
	return &Request{
		QueryParams: req.QueryStringParameters,
		Body:        req.Body,
		Headers:     req.Headers,
	}
}

func BuildAPIGWResponse(res Response, err error) (events.APIGatewayProxyResponse, error) {
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	body, err := json.Marshal(res.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: res.Status,
		Body:       string(body),
		Headers:    defaultHeaders(),
	}, nil
}

func defaultHeaders() map[string]string {
	return map[string]string{"Content-Type": "application/json"}
}
