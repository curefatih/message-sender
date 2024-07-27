package main

import (
	"context"

	"github.com/curefatih/message-sender/cache"
	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/handler"
	"github.com/curefatih/message-sender/model"
	"github.com/curefatih/message-sender/runner"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	cfg := readConfig()

	dbConn := db.InitPostgreSQLConnection(ctx, cfg.GetString("db.postgresql.dsn"))
	taskStateRepository := db.NewPostgreSQLTaskStateRepository(cfg, dbConn)
	messageTaskRepository := db.NewPostgreSQLMessageTaskRepository(cfg, dbConn)

	redisClient := cache.SetupRedisClient(ctx, cfg)
	messageTaskResultCache := cache.NewRedisCache[model.MessageTaskResult](redisClient, cfg)

	r := runner.NewSentMessageTaskRunner(ctx, cfg, messageTaskRepository, taskStateRepository, messageTaskResultCache)
	r.Run(ctx)

	router := gin.New()

	if err := handler.Setup(
		ctx,
		cfg,
		router,
		messageTaskRepository,
		taskStateRepository,
		messageTaskResultCache,
	).Run(); err != nil {
		log.Fatal().Err(err)
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
		log.Fatal().Msgf("fatal error config file: %w", err)
	}
	return viper.GetViper()
}
