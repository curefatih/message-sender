package model

import "gorm.io/gorm"

type MessageTask struct {
	gorm.Model
	MessageContent string
	PhoneNumber    string
}
