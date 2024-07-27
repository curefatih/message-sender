package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Middleware interface {
	AuthInsiderMiddleware(*viper.Viper) gin.HandlerFunc
}
