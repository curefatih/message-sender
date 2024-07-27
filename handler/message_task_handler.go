package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/curefatih/message-sender/cache"
	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/model"
	"github.com/curefatih/message-sender/model/dto"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type MessageTaskHandler struct {
	ctx                    context.Context
	cfg                    *viper.Viper
	repository             db.MessageTaskRepository
	messageTaskResultCache cache.Cache[model.MessageTaskResult]
}

func NewMessageTaskHandler(ctx context.Context, cfg *viper.Viper, repository db.MessageTaskRepository, messageTaskResultCache cache.Cache[model.MessageTaskResult]) *MessageTaskHandler {
	return &MessageTaskHandler{
		ctx:                    ctx,
		cfg:                    cfg,
		repository:             repository,
		messageTaskResultCache: messageTaskResultCache,
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
// @Param   MessageTaskCreateRequest body dto.MessageTaskCreateRequest true "MessageTaskCreateRequest"
// @Failure   400
// @Failure   500
// @Success 201 {object} dto.MessageTaskCreateResponse
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

	ctx.JSON(http.StatusCreated, dto.MessageTaskCreateResponse{
		Message: "successfully created new message task",
		ID:      res.ID,
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
// @Param			id	path	int	true	"message task id"
// @Success 	200 {string} OK
// @Failure   500
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
// @Summary Gets message tasks with pagination
// @Schemes
// @Description Gets message tasks with pagination
// @Tags Message Task
// @Accept json
// @Produce json
// @Param 	status 				query  string	true  	"status filter" Enums(COMPLETED, WAITING, PROCESSING, FAILED)
// @Param 	page 					query  number			true  	"page"  minimum(1)
// @Param 	page_size 		query  number  		true  	"page size" minimum(1)    maximum(20)
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

// @BasePath /api/v1
// DeleteMessageTask godoc
// @Summary Gets message task result from cache
// @Schemes
// @Description Gets message task result from cache
// @Tags Message Task
// @Accept json
// @Produce json
// @Param id path  string  true  "message id"
// @Success 200 {object} model.MessageTaskResult
// @Router /api/v1/tasks/messages/{id} [get]
func (mth *MessageTaskHandler) GetMessageTaskResult(ctx *gin.Context) {
	messageTaskID := ctx.Param("id")

	res, err := mth.messageTaskResultCache.Get(ctx.Request.Context(), messageTaskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "result couldn't find in cache",
		})
		return
	}

	fmt.Println("got res", res)

	ctx.JSON(http.StatusOK, model.MessageTaskResult{
		TaskID:      messageTaskID,
		MessageID:   res.MessageID,
		SendingTime: res.SendingTime,
	})
}
