package model

import "gorm.io/gorm"

type MessageTask struct {
	gorm.Model
	MessageContent string     `json:"content"`
	To             string     `json:"to"`
	Status         TaskStatus `json:"status", sql:"type:task_status"`
}
