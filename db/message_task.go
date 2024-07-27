package db

import (
	"context"

	"github.com/curefatih/message-sender/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type MessageTaskRepository interface {
	Create(ctx context.Context, messageTask *model.MessageTask) (*model.MessageTask, error)
	GetUnprocessedNMessageTaskAndMarkAsProcessing(ctx context.Context, n int) ([]*model.MessageTask, error)
	DeleteById(ctx context.Context, id string) error
}

type PostgreSQLMessageTaskRepository struct {
	db  *gorm.DB
	cfg *viper.Viper
}

func NewPostgreSQLMessageTaskRepository(cfg *viper.Viper, db *gorm.DB) *PostgreSQLMessageTaskRepository {
	return &PostgreSQLMessageTaskRepository{
		cfg: cfg,
		db:  db,
	}
}

var _ MessageTaskRepository = &PostgreSQLMessageTaskRepository{}

// Create implements MessageTaskRepository.
func (p *PostgreSQLMessageTaskRepository) Create(ctx context.Context, messageTask *model.MessageTask) (*model.MessageTask, error) {
	res := p.db.Create(messageTask)

	if res.Error != nil {
		return nil, res.Error
	}

	return messageTask, nil
}

// DeleteById implements MessageTaskRepository.
func (p *PostgreSQLMessageTaskRepository) DeleteById(ctx context.Context, id string) error {
	res := p.db.Delete(&model.MessageTask{}, id)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

// GetUnprocessedNMessageTaskAndMarkAsProcessing implements MessageTaskRepository.
func (p *PostgreSQLMessageTaskRepository) GetUnprocessedNMessageTaskAndMarkAsProcessing(ctx context.Context, n int) ([]*model.MessageTask, error) {
	panic("unimplemented")
}
