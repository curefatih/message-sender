package handler

import (
	"context"
	"net/http"

	"github.com/curefatih/message-sender/cmd/api/docs"
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
) *gin.Engine {
	// init swagger
	initSwagger(router)

	// Protect with header auth key
	authorized := router.Group("/")
	authorized.Use(m.AuthInsiderMiddleware(cfg))
	authorized.POST("/")
	authorized.DELETE("/")

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
	v1 := r.Group("/")
	{
		eg := v1.Group("/")
		{
			eg.GET("/health", Health)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
