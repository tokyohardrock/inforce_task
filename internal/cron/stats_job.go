package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	jobTrunc    = time.Hour
	jobInterval = 4 * jobTrunc
)

type StatsRepo interface {
	CalculateAndSaveUserStats(ctx context.Context, from, to time.Time) error
}

type StatsJob struct {
	cron *cron.Cron
	repo StatsRepo
}

func NewStatsJob(repo StatsRepo) *StatsJob {
	return &StatsJob{
		cron: cron.New(),
		repo: repo,
	}
}

func (j *StatsJob) Start() error {
	_, err := j.cron.AddFunc("0 */4 * * *", func() {
		j.run()
	})
	if err != nil {
		return err
	}

	j.cron.Start()
	slog.Info("cron scheduler started: running every 4 hours")
	return nil
}

func (j *StatsJob) Stop() {
	ctx := j.cron.Stop()
	<-ctx.Done()
	slog.Info("cron scheduler stopped gracefully")
}

func (j *StatsJob) run() {
	slog.Info("starting user activity stats calculation job")

	now := time.Now().UTC().Truncate(jobTrunc)
	from := now.Add(-jobInterval)
	to := now

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := j.repo.CalculateAndSaveUserStats(ctx, from, to)
	if err != nil {
		slog.Error("failed to calculate and save user stats", "error", err, "from", from, "to", to)
		return
	}

	slog.Info("user activity stats successfully calculated", "time_bucket", from)
}
