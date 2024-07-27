package dto

import "github.com/curefatih/message-sender/model"

type MessageTaskCreateRequest struct {
	MessageContent string `json:"content"`
	To             string `json:"to"`
}

func (mtcr *MessageTaskCreateRequest) ToMessageTask() *model.MessageTask {
	return &model.MessageTask{
		MessageContent: mtcr.MessageContent,
		To:             mtcr.To,
	}
}

type MessageTaskSendPayload struct {
	MessageContent string `json:"content"`
	To             string `json:"to"`
}
