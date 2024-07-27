package model

import "database/sql/driver"

type TaskStatus string

const (
	WAITING    TaskStatus = "WAITING"
	PROCESSING TaskStatus = "PROCESSING"
	COMPLETE   TaskStatus = "COMPLETE"
	FAILED     TaskStatus = "FAILED"
)

func (self *TaskStatus) Scan(value interface{}) error {
	*self = TaskStatus(value.([]byte))
	return nil
}

func (self TaskStatus) Value() (driver.Value, error) {
	return string(self), nil
}
