package domain

type TrackingInfo struct {
	Code           string `json:"code"`
	LastSearchDate string `json:"last_search_date"`
	LastEventDate  string `json:"last_event_date"`
	Users          []User `json:"users"`
	LastEvent      Event  `json:"last_event"`
}

type User struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

type Event struct {
	Status string `json:"status"`
	Place  string `json:"place"`
	Date   string `json:"date"`
	Time   string `json:"time"`
}
