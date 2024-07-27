package db

import (
	"context"

	"github.com/curefatih/message-sender/model"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

func InitPostgreSQLConnection(ctx context.Context, dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.MessageTask{})
	db.AutoMigrate(&model.TaskState{})

	return db
}
