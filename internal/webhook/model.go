package webhook

type WebhookRequestDTO struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	Id      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Field string `json:"field"`
	Value Value  `json:"value"`
}

type Value struct {
	MessagingProduct string            `json:"messaging_product"`
	Metadata         MetadataWebhook   `json:"metadata"`
	Contacts         *[]ContactWebhook `json:"contacts,omitempty"`
	Messages         *[]MessageWebhook `json:"messages,omitempty"`
	Statuses         *[]StatusWebhook  `json:"statuses,omitempty"`
}

type ContactWebhook struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type MessageWebhook struct {
	Button    *ButtonWebhook  `json:"button,omitempty"`
	Text      *TextWebhook    `json:"text,omitempty"`
	Context   *ContextWebhook `json:"context,omitempty"`
	From      string          `json:"from"`
	Id        string          `json:"id"`
	Timestamp string          `json:"timestamp"`
	Type      string          `json:"type"`
}

type TextWebhook struct {
	Body string `json:"body"`
}

type ButtonWebhook struct {
	Payload string `json:"payload"`
	Text    string `json:"text"`
}

type ContextWebhook struct {
	From string `json:"from"`
	Id   string `json:"id"`
}

type MetadataWebhook struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberId      string `json:"phone_number_id"`
}

type StatusWebhook struct {
}
