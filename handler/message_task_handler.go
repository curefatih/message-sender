package handler

import (
	"context"
	"net/http"

	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/model"
	"github.com/curefatih/message-sender/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type MessageTaskHandler struct {
	ctx        context.Context
	cfg        *viper.Viper
	repository db.MessageTaskRepository
}

func NewMessageTaskHandler(ctx context.Context, cfg *viper.Viper, repository db.MessageTaskRepository) *MessageTaskHandler {
	return &MessageTaskHandler{
		ctx:        ctx,
		cfg:        cfg,
		repository: repository,
	}
}

// @BasePath /api/v1
// CreateMessageTask godoc
// @Summary CreatesMessageTask
// @Schemes
// @Description Creates new message task that will be consumed.
// @Tags Message Task
// @Accept json
// @Produce json
// @Param   x-ins-auth-key header string true "Auth Key"
// @Param   MessageTaskCreateRequest body dto.MessageTaskCreateRequest true "Add MessageTaskCreateRequest"
// @Success 201 {object} model.MessageTask
// @Router /api/v1/tasks/messages/ [post]
func (mth *MessageTaskHandler) CreateMessageTask(ctx *gin.Context) {
	var messageTaskReq dto.MessageTaskCreateRequest

	err := ctx.ShouldBind(&messageTaskReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	messageTask := messageTaskReq.ToMessageTask()
	messageTask.Status = model.WAITING
	res, err := mth.repository.Create(ctx.Request.Context(), messageTask)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while creating new message task",
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "successfully created new message task",
		"id":      res.ID,
	})
}

// @BasePath /api/v1
// DeleteMessageTask godoc
// @Summary Deletes Message Task
// @Schemes
// @Description Deletes message task that will
// @Tags Message Task
// @Accept json
// @Produce json
// @Param   x-ins-auth-key header string true "Auth Key"
// @Param        id   		path      int  true  "Message Task ID"
// @Success 200 {string} OK
// @Router /api/v1/tasks/messages/{id} [delete]
func (mth *MessageTaskHandler) DeleteMessageTask(ctx *gin.Context) {

	messageTaskID := ctx.Param("id")

	err := mth.repository.DeleteById(ctx.Request.Context(), messageTaskID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while deleting message task",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "message task successfully deleted",
	})
}
