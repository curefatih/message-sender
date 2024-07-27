package model

import "database/sql/driver"

type TaskStatus string

const (
	TaskStatusWaiting    TaskStatus = "WAITING"
	TaskStatusProcessing TaskStatus = "PROCESSING"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusFailed     TaskStatus = "FAILED"
)

func (self *TaskStatus) Scan(value interface{}) error {
	*self = TaskStatus(value.([]byte))
	return nil
}

func (self TaskStatus) Value() (driver.Value, error) {
	return string(self), nil
}
