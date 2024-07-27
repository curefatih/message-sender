package model

import "gorm.io/gorm"

type MessageTask struct {
	gorm.Model
	MessageContent string `json:"content"`
	To             string `json:"to"`
}
