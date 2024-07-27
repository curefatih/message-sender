package db

import (
	"context"
	"time"

	"github.com/curefatih/message-sender/model"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TaskStateRepository interface {
	UpdateTaskActiveStatus(ctx context.Context, active bool) error
	GetOrCreateTaskState(ctx context.Context) (*model.TaskState, error)
	UpdateTaskState(ctx context.Context, taskState *model.TaskState) error
}

type PostgreSQLTaskStateRepository struct {
	db  *gorm.DB
	cfg *viper.Viper
}

var _ TaskStateRepository = &PostgreSQLTaskStateRepository{}

func NewPostgreSQLTaskStateRepository(cfg *viper.Viper, db *gorm.DB) *PostgreSQLTaskStateRepository {
	return &PostgreSQLTaskStateRepository{
		cfg: cfg,
		db:  db,
	}
}

// UpdateTaskActiveStatus implements TaskStateRepository.
func (p *PostgreSQLTaskStateRepository) UpdateTaskActiveStatus(ctx context.Context, active bool) error {
	taskState, err := p.GetOrCreateTaskState(ctx)
	if err != nil {
		return err
	}

	taskState.Active = active
	res := p.db.Save(taskState)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// GetOrCreateTaskState implements TaskStateRepository.
func (p *PostgreSQLTaskStateRepository) GetOrCreateTaskState(ctx context.Context) (*model.TaskState, error) {
	var messageTask model.TaskState

	result := p.db.First(&messageTask)

	// If no record is found, create a new one
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			messageTask = model.TaskState{
				Active:                  true,
				LastSuccessfulQueryTime: time.Now(),
				Status:                  model.WAITING,
			}
			createResult := p.db.Create(&messageTask)
			if createResult.Error != nil {
				log.Fatal().Msgf("failed to create user: %v", createResult.Error)
			}

			return &messageTask, nil
		} else {
			log.Fatal().Msgf("failed to fetch user: %v", result.Error)
		}
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &messageTask, nil
}

// UpdateTaskStateWithTx implements TaskStateRepository.
func (p *PostgreSQLTaskStateRepository) UpdateTaskState(ctx context.Context, taskState *model.TaskState) error {

	currentTaskState, err := p.GetOrCreateTaskState(ctx)
	if err != nil {
		log.Info().Msg("error while getting task state")
		return err
	}

	taskState.ID = currentTaskState.ID

	res := p.db.Updates(taskState)
	if res.Error != nil {
		log.Info().Msg("error while updating task state")
		return res.Error
	}

	return nil
}
