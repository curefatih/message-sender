package model

type CreateMessageRequestPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type CreateMessageResponsePayload struct {
	Message   string `json:"message"`
	MessageId string `json:"messageId"`
}
