package main

import (
	"context"
	"fmt"
	"log"

	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/handler"
	"github.com/curefatih/message-sender/runner"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	cfg := readConfig()

	dbConn := db.InitPostgreSQLConnection(ctx, cfg.GetString("db.postgresql.dsn"))
	taskStateRepository := db.NewPostgreSQLTaskStateRepository(cfg, dbConn)
	messageTaskRepository := db.NewPostgreSQLMessageTaskRepository(cfg, dbConn)

	r := runner.NewSentMessageTaskRunner(ctx, cfg, messageTaskRepository, taskStateRepository)
	r.Run(ctx)

	router := gin.Default()

	if err := handler.Setup(
		ctx,
		cfg,
		router,
		messageTaskRepository,
		taskStateRepository,
	).Run(); err != nil {
		log.Fatal(err)
		r.Stop()
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
