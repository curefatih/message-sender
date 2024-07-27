package model

import "time"

type TaskState struct {
	LastSuccessfulQueryTime time.Time
	Status                  TaskStatus
}
