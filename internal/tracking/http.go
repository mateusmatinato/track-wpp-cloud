package tracking

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client interface {
	Track(ctx context.Context, code string) (TrackClientResponse, error)
}

const (
	trackUrl = "https://api.linketrack.com/track/json?user=%s&token=%s&codigo=%s"
)

type client struct {
	Properties ClientProperties
}

func (c client) Track(ctx context.Context, code string) (TrackClientResponse, error) {
	log.Printf("[INFO] starting track client | code: %s", code)

	resp, err := http.Get(fmt.Sprintf(trackUrl, c.Properties.Username, c.Properties.Password, code))
	if err != nil {
		log.Printf("[ERROR] Error track client - http | error: %s", err.Error())
		return TrackClientResponse{}, err
	}
	defer resp.Body.Close()

	responseJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] Error track client - read body | error: %s", err.Error())
		return TrackClientResponse{}, err
	}

	var response TrackClientResponse
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		log.Printf("[ERROR] Error track client - unmarshal | error: %s", err.Error())
		return TrackClientResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] Error track client - status code | response_code: %s", resp.StatusCode)
		return TrackClientResponse{}, errors.New("error track client")
	}

	return response, nil
}

func NewClient() Client {
	return &client{Properties: ClientProperties{
		Username: "USER",
		Password: "PWD",
	}}
}
