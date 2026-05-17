package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"

	"github.com/feishu/feishu-office-suite/internal/data"
	"github.com/feishu/feishu-office-suite/internal/job"
)

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path")
}

func main() {
	flag.Parse()

	logger := log.NewStdLogger(os.Stdout)
	logHelper := log.NewHelper(logger)

	cfg, err := data.NewConfig(flagconf)
	if err != nil {
		logHelper.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	dataData, cleanup, err := data.NewData(logHelper, cfg, nil)
	if err != nil {
		logHelper.Errorf("Failed to create data: %v", err)
		os.Exit(1)
	}
	defer cleanup()

	executor := job.NewExecutor(dataData)

	c := cron.New()
	_, err = c.AddFunc("0 */5 * * * *", func() {
		logHelper.Info("Running scheduled user sync job")
		task := &job.Task{
			Type:    job.TaskTypeSyncUser,
			Payload: job.TaskPayload{},
		}
		if err := executor.ExecuteTask(context.Background(), task); err != nil {
			logHelper.Errorf("Failed to sync users: %v", err)
		}
	})
	if err != nil {
		logHelper.Errorf("Failed to add cron job: %v", err)
	}

	_, err = c.AddFunc("0 */10 * * * *", func() {
		logHelper.Info("Running scheduled department sync job")
		task := &job.Task{
			Type:    job.TaskTypeSyncDepartment,
			Payload: job.TaskPayload{},
		}
		if err := executor.ExecuteTask(context.Background(), task); err != nil {
			logHelper.Errorf("Failed to sync departments: %v", err)
		}
	})
	if err != nil {
		logHelper.Errorf("Failed to add cron job: %v", err)
	}

	c.Start()
	logHelper.Info("Worker started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logHelper.Info("Shutting down worker...")
	c.Stop()
}