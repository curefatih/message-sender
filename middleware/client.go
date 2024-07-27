package middleware

import (
	"net/http"

	"github.com/curefatih/message-sender/handler"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type ClientMiddleware struct{}

var _ handler.Middleware = &ClientMiddleware{}

// AuthInsiderMiddleware implements handler.Middleware.
func (c ClientMiddleware) AuthInsiderMiddleware(cfg *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authKey := ctx.GetHeader("x-ins-auth-key")
		if authKey != cfg.GetString("auth.key") {
			ctx.JSON(http.StatusForbidden, "Not authorized")
			ctx.Abort()
		}
		ctx.Next()
	}
}
