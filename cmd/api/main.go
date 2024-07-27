package main

import (
	"context"
	"fmt"
	"log"

	"github.com/curefatih/message-sender/handler"
	"github.com/curefatih/message-sender/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	cfg := readConfig()

	router := gin.Default()
	m := middleware.ClientMiddleware{}

	if err := handler.Setup(
		ctx,
		cfg,
		router,
		m,
	).Run(); err != nil {
		log.Fatal(err)
		return
	}
}

func readConfig() *viper.Viper {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Fatal(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}
