package tracking

import (
	"context"
	"errors"
	"log"
)

type Service interface {
	TrackPackage(ctx context.Context, code string) (TrackPackageResult, error)
}

type service struct {
	Client Client
}

var (
	errTrackEmpty = errors.New("error - package doesnt have any event yet")
)

func (s service) TrackPackage(ctx context.Context, code string) (TrackPackageResult, error) {
	log.Printf("[INFO] starting track service | code: %s", code)

	clientResponse, err := s.Client.Track(ctx, code)
	if err != nil {
		return TrackPackageResult{}, err
	}

	if len(clientResponse.Events) == 0 {
		return TrackPackageResult{}, errTrackEmpty
	}

	lastEvent := clientResponse.Events[0]

	result := TrackPackageResult{
		LastUpdate: clientResponse.LastUpdate,
		Event:      lastEvent,
	}

	return result, nil
}

func NewService(client Client) Service {
	return &service{Client: client}
}
