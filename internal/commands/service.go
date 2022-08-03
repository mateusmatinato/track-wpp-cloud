package commands

import (
	"context"
	"log"
	"strings"
	"time"
	"wpp-cloud/internal/domain"
	"wpp-cloud/internal/message"
	"wpp-cloud/internal/repository"
	"wpp-cloud/internal/tracking"
)

var commandsAvailable = map[string]Command{
	"/rastrear": RegisterTrack,
	"/ajuda":    Help,
}

type Command int64

const (
	RegisterTrack Command = iota
	Help
)

type Service interface {
	ProcessMessage(ctx context.Context, request ProcessRequest) error
}

type service struct {
	TrackerService tracking.Service
	MessageService message.Service
	Repository     repository.Repository
}

func (s service) ProcessMessage(ctx context.Context, request ProcessRequest) error {

	args := strings.Split(request.Message, " ")
	if len(args) == 0 {
		log.Printf("[ERROR] error args empty")
		return s.sendGenericErrorMessage(ctx, request.Contact)
	}

	command, exist := commandsAvailable[args[0]]
	if !exist {
		log.Printf("[ERROR] error command unknown")
		return s.sendGenericErrorMessage(ctx, request.Contact)
	}

	log.Printf("[INFO] starting process command | command: %s | request: %s", args[0], request)
	switch command {
	case RegisterTrack:
		return s.registerTrack(ctx, request.Contact, args...)
	case Help:
		return s.help(ctx, request.Contact)
	}

	return nil
}

func (s service) help(ctx context.Context, contact Contact) error {
	return s.MessageService.SendMessage(ctx, contact.Number, message.GetHelp())
}

func (s service) registerTrack(ctx context.Context, contact Contact, args ...string) error {
	if len(args) > 3 || len(args) < 2 {
		log.Printf("[ERROR] error args size | args: %s", args)
		err := s.MessageService.SendMessage(ctx, contact.Number, message.GetTrackSteps())
		if err != nil {
			return err
		}
		return nil
	}

	trackingCode, packageName := extractVariablesFromArgs(args)

	trackingInfo, err := s.Repository.GetTrackingInfo(trackingCode)
	if err != nil {
		return s.MessageService.SendMessage(ctx, contact.Number, message.GetGenericTrackError(trackingCode))
	}

	log.Printf("[INFO] starting register track | code: %s | name: %s", trackingCode, packageName)
	result, err := s.TrackerService.TrackPackage(ctx, trackingCode)
	if err != nil {
		return s.MessageService.SendMessage(ctx, contact.Number, message.GetGenericTrackError(trackingCode))
	}

	log.Printf("[INFO] Track success | result: %v", result)
	err = s.MessageService.SendMessage(ctx, contact.Number, message.GetTrackSuccess(trackingCode, result))
	if err != nil {
		return err
	}

	trackingInfo = domain.TrackingInfo{
		Code:           trackingCode,
		LastSearchDate: time.Now().String(),
		LastEventDate:  result.LastUpdate.String(),
		Users:          generateUsers(trackingInfo.Users, contact),
		LastEvent: domain.Event{
			Status: result.Event.Status,
			Place:  result.Event.Place,
			Date:   result.Event.Date,
			Time:   result.Event.Time,
		},
	}

	err = s.Repository.SaveTrackingInfo(trackingInfo)
	if err != nil {
		return err
	}

	return nil
}

func generateUsers(users []domain.User, contact Contact) []domain.User {
	found := false
	for _, user := range users {
		if user.Number == contact.Number {
			found = true
		}
	}

	if !found {
		return append(users, domain.User{
			Name:   contact.Name,
			Number: contact.Number,
		})
	}
	return users
}

func extractVariablesFromArgs(args []string) (string, string) {
	trackingCode := args[1]
	var packageName string
	if len(args) == 3 {
		packageName = args[2]
	} else {
		packageName = trackingCode
	}
	return trackingCode, packageName
}

func (s service) validateArguments(ctx context.Context, contact Contact, args []string) bool {

	return true
}

func (s service) sendGenericErrorMessage(ctx context.Context, contact Contact) error {
	return s.MessageService.SendMessage(ctx, contact.Number, message.GetGenericMessageError())
}

func NewService(
	trackerService tracking.Service,
	messageService message.Service,
	repository repository.Repository,
) Service {
	return &service{
		TrackerService: trackerService,
		MessageService: messageService,
		Repository:     repository,
	}
}
