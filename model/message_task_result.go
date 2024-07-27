package model

import "time"

type MessageTaskResult struct {
	TaskID      string    `json:"taskId"`
	MessageID   string    `json:"messageId"`
	SendingTime time.Time `json:"sendingTime"`
}
