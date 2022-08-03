package message

type WppMessageRequest struct {
	MessagingProduct string `json:"messaging_product"`
	To               string `json:"to"`
	Text             Text   `json:"text"`
}

type Text struct {
	Body string `json:"body"`
}

type WppMessageResponse struct {
	MessagingProduct string    `json:"messaging_product"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Contact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type Message struct {
	Id string `json:"id"`
}
