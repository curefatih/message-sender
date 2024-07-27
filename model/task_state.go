package model

import (
	"time"

	"gorm.io/gorm"
)

type TaskState struct {
	gorm.Model
	LastSuccessfulQueryTime time.Time
	Status                  TaskStatus
	Active                  bool
}
