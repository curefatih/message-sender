package model

import (
	"database/sql/driver"
)

type TaskStatus string

const (
	WAITING    TaskStatus = "WAITING"
	PROCESSING TaskStatus = "PROCESSING"
	COMPLETED  TaskStatus = "COMPLETED"
	FAILED     TaskStatus = "FAILED"
)

func (self *TaskStatus) Scan(value interface{}) error {
	*self = TaskStatus(value.(string))
	return nil
}

func (self TaskStatus) Value() (driver.Value, error) {
	return string(self), nil
}
