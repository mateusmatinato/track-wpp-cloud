package recurring

import (
	"context"
	"fmt"
	"time"
	"wpp-cloud/internal/domain"
	"wpp-cloud/internal/message"
	"wpp-cloud/internal/repository"
	"wpp-cloud/internal/tracking"
)

type Service interface {
	ProcessCodes(ctx context.Context) error
}

type service struct {
	TrackerService tracking.Service
	MessageService message.Service
	Repository     repository.Repository
}

func (s service) ProcessCodes(ctx context.Context) error {
	fmt.Printf("[INFO] Starting recurring service\n")

	infos, err := s.Repository.ScanTrackingInfo()
	if err != nil {
		return err
	}

	if len(infos) == 0 {
		fmt.Printf("[INFO] We dont have any codes to track\n")
		return nil
	}

	fmt.Printf("[INFO] Starting to track codes | amount: %d\n", len(infos))
	for _, info := range infos {
		result, err := s.TrackerService.TrackPackage(ctx, info.Code)
		if err != nil {
			break
		}
		fmt.Printf("[INFO] Tracked package | code: %s | tracking: %v\n", info.Code, result)
		if result.LastUpdate.String() != info.LastEventDate {
			fmt.Printf("[INFO] Has update, starting notificate users | code: %s", info.Code)
			for _, user := range info.Users {
				fmt.Printf("[INFO] Notificating user about package | code: %s | user: %s", info.Code, user.Name+"-"+user.Number)
				err = s.MessageService.SendMessage(ctx, user.Number, message.GetTrackUpdateSuccess(info.Code, result))
				if err != nil {
					break
				}

				info.LastSearchDate = time.Now().String()
				info.LastEventDate = result.LastUpdate.String()
				info.LastEvent = domain.Event{
					Status: result.Event.Status,
					Place:  result.Event.Place,
					Date:   result.Event.Date,
					Time:   result.Event.Time,
				}
				err := s.Repository.SaveTrackingInfo(info)
				if err != nil {
					fmt.Printf("[ERROR] Error updating information about package | code: %s | error: %s\n", info.Code, err.Error())
					break
				}
			}
		} else {
			fmt.Printf("[INFO] Doesnt has update, ignoring | code: %s", info.Code)
		}
	}
	return nil
}

func NewService(trackService tracking.Service, messageService message.Service, repo repository.Repository) Service {
	return &service{
		TrackerService: trackService,
		MessageService: messageService,
		Repository:     repo,
	}
}
