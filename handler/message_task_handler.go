package handler

import (
	"context"

	"github.com/curefatih/message-sender/db"
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

func (mth *MessageTaskHandler) CreateMessageTask(ctx *gin.Context) {}

func (mth *MessageTaskHandler) GetSentMessageTasks(ctx *gin.Context) {}

func (mth *MessageTaskHandler) DeleteMessageTask(ctx *gin.Context) {}
