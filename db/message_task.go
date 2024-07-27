package db

import (
	"context"
	"fmt"

	"github.com/curefatih/message-sender/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type MessageTaskRepository interface {
	Create(ctx context.Context, messageTask *model.MessageTask) (*model.MessageTask, error)
	GetUnprocessedNMessageTaskAndMarkAsProcessing(ctx context.Context, n int64) ([]*model.MessageTask, error)
	GetPaginated(ctx context.Context, page, pageSize int, status *string) ([]*model.MessageTask, error)
	DeleteById(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status model.TaskStatus) error
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
func (p *PostgreSQLMessageTaskRepository) GetUnprocessedNMessageTaskAndMarkAsProcessing(ctx context.Context, n int64) ([]*model.MessageTask, error) {
	var tasks []*model.MessageTask
	var err error

	// Start a transaction
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Ensure the transaction is rolled back in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Fetch tasks with PROCESSING status first
	err = tx.Where("status = ?", model.PROCESSING).
		Limit(int(n)).
		Find(&tasks).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch PROCESSING tasks: %w", err)
	}

	// If less than n tasks were fetched, fetch additional tasks with WAITING status
	if int64(len(tasks)) < n {
		remaining := n - int64(len(tasks))
		var additionalTasks []*model.MessageTask

		err = tx.Where("status = ?", model.WAITING).
			Limit(int(remaining)).
			Find(&additionalTasks).
			Error

		if err != nil {
			return nil, fmt.Errorf("failed to fetch WAITING tasks: %w", err)
		}

		tasks = append(tasks, additionalTasks...)
	}

	if len(tasks) == 0 {
		return tasks, nil
	}

	// Mark tasks as PROCESSING
	for _, task := range tasks {
		task.Status = model.PROCESSING
	}

	err = tx.Save(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tasks, nil
}

// UpdateStatus implements MessageTaskRepository.
func (p *PostgreSQLMessageTaskRepository) UpdateStatus(ctx context.Context, id string, status model.TaskStatus) error {
	// Use the context with the database query
	if err := p.db.WithContext(ctx).Model(&model.MessageTask{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// GetPaginated implements MessageTaskRepository.
func (p *PostgreSQLMessageTaskRepository) GetPaginated(ctx context.Context, page int, pageSize int, status *string) ([]*model.MessageTask, error) {
	offset := (page - 1) * pageSize

	var tasks []*model.MessageTask

	query := p.db.WithContext(ctx).Model(&model.MessageTask{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// Fetch paginated tasks
	err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&tasks).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch paginated tasks: %w", err)
	}

	return tasks, nil
}
