package db

import (
	"context"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type TaskStateRepository interface {
	UpdateTaskActiveStatus(ctx context.Context, active bool) error
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
	panic("unimplemented")
}
