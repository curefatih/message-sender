package runner

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/curefatih/message-sender/cache"
	"github.com/curefatih/message-sender/db"
	"github.com/curefatih/message-sender/model"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

type SentMessageTaskRunner struct {
	cron                   *cron.Cron
	ctx                    context.Context
	cfg                    *viper.Viper
	messageTaskRepository  db.MessageTaskRepository
	taskStateRepository    db.TaskStateRepository
	messageTaskResultCache cache.Cache[model.MessageTaskResult]
}

func NewSentMessageTaskRunner(
	ctx context.Context,
	cfg *viper.Viper,
	messageTaskRepository db.MessageTaskRepository,
	taskStateRepository db.TaskStateRepository,
	messageTaskResultCache cache.Cache[model.MessageTaskResult],
) *SentMessageTaskRunner {
	c := cron.New(cron.WithSeconds())
	return &SentMessageTaskRunner{
		cron:                   c,
		ctx:                    ctx,
		cfg:                    cfg,
		messageTaskRepository:  messageTaskRepository,
		taskStateRepository:    taskStateRepository,
		messageTaskResultCache: messageTaskResultCache,
	}
}

var _ Runner = &SentMessageTaskRunner{}

// Run implements Runner.
func (s *SentMessageTaskRunner) Run(ctx context.Context) error {
	_, err := s.cron.AddFunc(s.cfg.GetString("process.task.cron"), func() {
		log.Info().Msg("running sent message task")
		count := s.cfg.GetInt64("process.task.count")
		retryCount := s.cfg.GetInt("process.task.retry")
		deltaMin := s.cfg.GetInt64("process.task.delta_in_minutes")

		taskState, err := s.taskStateRepository.GetOrCreateTaskState(ctx)

		if err != nil {
			log.Fatal().Msgf("couldn't find task state. terminating...", err)
		}

		if !taskState.Active {
			log.Info().Msg("task schedule inactive. skipping...")
			return
		}

		if !shouldRunTask(taskState.LastSuccessfulQueryTime, deltaMin) {
			log.Info().Msg("delta not in past. tasks waiting next turn...")
			return
		}

		currentTime := time.Now()
		tasks, err := s.messageTaskRepository.GetUnprocessedNMessageTaskAndMarkAsProcessing(ctx, count)
		if err != nil {
			// we may not exit?
			log.Fatal().Msg("couldn't get message tasks. terminating...")
		}

		// Create a new errgroup
		g, ctx := errgroup.WithContext(context.Background())

		for _, task := range tasks {
			task := task // Create a new instance for the goroutine
			g.Go(func() error {
				err := invokeTask(ctx, s.cfg, task, retryCount, s.messageTaskResultCache)
				if err != nil {
					err := s.messageTaskRepository.UpdateStatus(ctx, strconv.FormatUint(uint64(task.ID), 10), model.FAILED)
					if err != nil {
						log.Fatal().Msg("couldn't update message task status. terminating...")
					}
					return nil
				}

				err = s.messageTaskRepository.UpdateStatus(ctx, strconv.FormatUint(uint64(task.ID), 10), model.COMPLETED)
				if err != nil {
					log.Fatal().Msg("couldn't update message task status. terminating...")
				}

				return nil
			})
		}

		err = s.taskStateRepository.UpdateTaskState(ctx, &model.TaskState{
			LastSuccessfulQueryTime: currentTime,
		})

		if err != nil {
			log.Info().Err(err)
			return
		}

		// Wait for all tasks to complete
		if err := g.Wait(); err != nil {
			log.Info().Err(err)
		} else {
			log.Info().Msg("all tasks completed successfully")
		}
	})

	if err != nil {
		return err
	}
	s.cron.Start()

	return nil
}

func invokeTask(
	ctx context.Context,
	cfg *viper.Viper,
	messageTask *model.MessageTask,
	retryCount int,
	messageTaskResultCache cache.Cache[model.MessageTaskResult],
) error {
	for i := 0; i <= retryCount; i++ {
		log.Info().Msgf("Invoking task %d, attempt %d\n", messageTask.ID, i+1)
		currentTime := time.Now()

		resp, err := runMessageTask(ctx, cfg, *messageTask)
		if err == nil {
			taskIDStr := strconv.FormatUint(uint64(messageTask.ID), 10)
			messageTaskResultCache.Set(ctx, taskIDStr, model.MessageTaskResult{
				TaskID:      taskIDStr,
				MessageID:   resp.MessageId,
				SendingTime: currentTime,
			})
			return nil
		}

		// If the context is done, return its error
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		log.Info().Msgf("Task %d failed on attempt %d: %v\n", messageTask.ID, i+1, err)
	}
	return fmt.Errorf("task %d failed after %d attempts", messageTask.ID, retryCount)
}

// Stop implements Runner.
func (s *SentMessageTaskRunner) Stop() error {
	s.cron.Stop()
	return nil
}

func shouldRunTask(scheduledTime time.Time, deltaMin int64) bool {
	// Add deltaMin minutes to the scheduled time
	deadline := scheduledTime.Add(time.Duration(deltaMin) * time.Minute)
	// Get the current time
	currentTime := time.Now()
	// Check if current time is past the deadline
	return currentTime.After(deadline)
}
