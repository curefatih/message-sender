package handler

import (
	"context"
	"net/http"

	"github.com/curefatih/message-sender/cmd/api/docs"
	"github.com/curefatih/message-sender/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(
	ctx context.Context,
	cfg *viper.Viper,
	router *gin.Engine,
	m Middleware,
	messageTaskRepository db.MessageTaskRepository,
	taskStateRepository db.TaskStateRepository,
) *gin.Engine {
	// init swagger

	r := router.Group("/api/v1")
	// health
	router.GET("/health", Health)

	taskStateHandler := NewTaskStateHandler(ctx, cfg, taskStateRepository)
	messageTaskHandler := NewMessageTaskHandler(ctx, cfg, messageTaskRepository)

	taskStateV1Endpoint := r.Group("/tasks")
	taskStateV1Endpoint.Use(m.AuthInsiderMiddleware(cfg))
	taskStateV1Endpoint.PUT("/", taskStateHandler.UpdateTaskState)

	// Protect with header auth key
	messageTasksV1Endpoint := taskStateV1Endpoint.Group("/messages")
	messageTasksV1Endpoint.POST("/", messageTaskHandler.CreateMessageTask)
	messageTasksV1Endpoint.DELETE("/:id", messageTaskHandler.DeleteMessageTask)

	// Task state

	initSwagger(router)

	return router
}

// @BasePath /api/v1

// PingExample godoc
// @Summary health
// @Schemes
// @Description response ok
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} health
// @Router /health [get]
func Health(g *gin.Context) {
	g.JSON(http.StatusOK, "OK")
}

func initSwagger(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
