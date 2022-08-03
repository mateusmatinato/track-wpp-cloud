package tracking

import "time"

type TrackPackageResult struct {
	LastUpdate time.Time
	Event
}

type Event struct {
	Status string `json:"status"`
	Place  string `json:"local"`
	Date   string `json:"data"`
	Time   string `json:"hora"`
}

type TrackClientResponse struct {
	Code        string    `json:"codigo"`
	Host        string    `json:"host"`
	Events      []Event   `json:"eventos"`
	Time        float64   `json:"time"`
	EventAmount int64     `json:"quantidade"`
	Service     string    `json:"servico"`
	LastUpdate  time.Time `json:"ultimo"`
}

type ClientProperties struct {
	Username string
	Password string
}
