package handler

import (
	"context"
	"net/http"

	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/model/dto"
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
// @Param   TaskStateUpdateRequest body dto.TaskStateUpdateRequest true "Add TaskStateUpdateRequest"
// @Success 200 {object} dto.TaskStateUpdateResponse
// @Router /api/v1/tasks [put]
func (tsh *TaskStateHandler) UpdateTaskState(ctx *gin.Context) {
	var taskStateUpdateRequest dto.TaskStateUpdateRequest

	err := ctx.ShouldBind(&taskStateUpdateRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = tsh.repository.UpdateTaskActiveStatus(ctx.Request.Context(), taskStateUpdateRequest.Active)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while creating new message task",
		})
	}

	ctx.JSON(http.StatusOK, dto.TaskStateUpdateResponse{
		Message: "successfully updated task state active status",
		Active:  taskStateUpdateRequest.Active,
	})
}
