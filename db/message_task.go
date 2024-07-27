package db

import (
	"context"

	"github.com/curefatih/message-sender/model"
)

type MessageTaskRepository interface {
	Create(ctx context.Context, trail *model.MessageTask) (*model.MessageTask, error)
	GetById(ctx context.Context, id string) (*model.MessageTask, error)
	UpdateById(ctx context.Context, trail *model.MessageTask) (*model.MessageTask, error)
	DeleteById(ctx context.Context, id string) error
}
