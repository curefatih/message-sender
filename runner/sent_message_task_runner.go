package runner

import (
	"context"

	"github.com/curefatih/message-sender/db"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
)

type SentMessageTaskRunner struct {
	cron                  *cron.Cron
	ctx                   context.Context
	cfg                   *viper.Viper
	messageTaskRepository db.MessageTaskRepository
	taskStateRepository   db.TaskStateRepository
}

func NewSentMessageTaskRunner(
	ctx context.Context,
	cfg *viper.Viper,
	messageTaskRepository db.MessageTaskRepository,
	taskStateRepository db.TaskStateRepository,
) *SentMessageTaskRunner {
	c := cron.New()
	return &SentMessageTaskRunner{
		cron:                  c,
		ctx:                   ctx,
		cfg:                   cfg,
		messageTaskRepository: messageTaskRepository,
		taskStateRepository:   taskStateRepository,
	}
}

var _ Runner = &SentMessageTaskRunner{}

// Run implements Runner.
func (s *SentMessageTaskRunner) Run(ctx context.Context) error {
	_, err := s.cron.AddFunc(s.cfg.GetString("process.task.cron"), func() {
		println("TODO: running now")
	})

	if err != nil {
		return err
	}
	s.cron.Start()

	return nil
}

// Stop implements Runner.
func (s *SentMessageTaskRunner) Stop() error {
	s.cron.Stop()
	return nil
}
