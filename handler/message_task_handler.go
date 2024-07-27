package handler

import (
	"context"
	"fmt"
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
// @Param   MessageTaskCreateRequest body dto.MessageTaskCreateRequest true "Add MessageTaskCreateRequest"
// @Success 201 {object} model.MessageTask
// @Router /api/v1/tasks/messages/ [post]
func (mth *MessageTaskHandler) CreateMessageTask(ctx *gin.Context) {
	maxContentLength := mth.cfg.GetInt("process.task.message.max_content_length")
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

	if len(messageTask.MessageContent) > maxContentLength {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("content length bigger than %d", maxContentLength),
		})
		return
	}

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

// @BasePath /api/v1
// DeleteMessageTask godoc
// @Summary Deletes Message Task
// @Schemes
// @Description Deletes message task that will
// @Tags Message Task
// @Accept json
// @Produce json
// @Param status query  object  true  "COMPLETED"
// @Success 200 {object} dto.PageResponse[model.MessageTask]
// @Router /api/v1/tasks/messages [get]
func (mth *MessageTaskHandler) GetMessagesWithPagination(ctx *gin.Context) {
	pageQuery := getPageQuery(ctx)
	status := ctx.Query("status")

	res, err := mth.repository.GetPaginated(ctx.Request.Context(), pageQuery.Page, pageQuery.PageSize, &status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while getting message tasks",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.PageResponse[model.MessageTask]{
		Page:     pageQuery.Page,
		PageSize: pageQuery.PageSize,
		Data:     res,
	})
}
