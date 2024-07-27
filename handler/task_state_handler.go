package handler

import (
	"context"

	"github.com/curefatih/message-sender/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type TaskStateHandler struct {
	ctx        context.Context
	cfg        *viper.Viper
	repository db.TaskStateRepository
}

func NewTaskStateHandler(ctx context.Context, cfg *viper.Viper, repository db.TaskStateRepository) *TaskStateHandler {
	return &TaskStateHandler{
		ctx:        ctx,
		cfg:        cfg,
		repository: repository,
	}
}

// @BasePath /api/v1

// UpdateTaskState godoc
// @Summary Updates Task State
// @Schemes
// @Description Creates new message task that will be consumed.
// @Tags Task State
// @Accept json
// @Produce json
// @Success 200 {string} todo
// @Router /api/v1/tasks [put]
func (mth *TaskStateHandler) UpdateTaskState(ctx *gin.Context) {
	panic("todo")
}
