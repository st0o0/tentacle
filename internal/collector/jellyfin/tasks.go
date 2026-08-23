package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type TasksCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	state        *prometheus.Desc
	progress     *prometheus.Desc
	lastDuration *prometheus.Desc
	lastSuccess  *prometheus.Desc
	lastRun      *prometheus.Desc
}

func NewTasksCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *TasksCollector {
	return &TasksCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		state: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "state"),
			"State of a scheduled task (0=Idle, 1=Running, 2=Cancelling).",
			[]string{"task_name", "category"}, nil,
		),
		progress: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "progress_ratio"),
			"Current progress of a running scheduled task (0 to 1).",
			[]string{"task_name"}, nil,
		),
		lastDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "last_run_duration_seconds"),
			"Duration of the last execution of a scheduled task.",
			[]string{"task_name"}, nil,
		),
		lastSuccess: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "last_run_success"),
			"Whether the last execution of a scheduled task was successful.",
			[]string{"task_name"}, nil,
		),
		lastRun: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "last_run_timestamp_seconds"),
			"Unix timestamp of when a scheduled task last completed.",
			[]string{"task_name"}, nil,
		),
	}
}

func (c *TasksCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.state
	ch <- c.progress
	ch <- c.lastDuration
	ch <- c.lastSuccess
	ch <- c.lastRun
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *TasksCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	tasks, err := c.client.GetScheduledTasks(ctx)
	duration := time.Since(start).Seconds()

	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "tasks")

	if err != nil {
		c.logger.Error("tasks collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "tasks")
		return
	}

	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "tasks")

	for _, t := range tasks {
		stateVal := 0.0
		switch t.State {
		case "Running":
			stateVal = 1.0
		case "Cancelling":
			stateVal = 2.0
		}
		ch <- prometheus.MustNewConstMetric(c.state, prometheus.GaugeValue, stateVal, t.Name, t.Category)
		ch <- prometheus.MustNewConstMetric(c.progress, prometheus.GaugeValue, t.CurrentProgressPercentage/100.0, t.Name)

		if t.LastExecutionResult != nil {
			successVal := 0.0
			if t.LastExecutionResult.Status == "Completed" {
				successVal = 1.0
			}
			ch <- prometheus.MustNewConstMetric(c.lastSuccess, prometheus.GaugeValue, successVal, t.Name)

			if endTime, ok := parseJellyfinTime(t.LastExecutionResult.EndTimeUtc); ok {
				ch <- prometheus.MustNewConstMetric(c.lastRun, prometheus.GaugeValue, float64(endTime.Unix()), t.Name)

				if startTime, ok := parseJellyfinTime(t.LastExecutionResult.StartTimeUtc); ok {
					ch <- prometheus.MustNewConstMetric(c.lastDuration, prometheus.GaugeValue, endTime.Sub(startTime).Seconds(), t.Name)
				}
			}
		}
	}
}
