package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type UsersCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	usersTotal   *prometheus.Desc
	usersOnline  *prometheus.Desc
	lastSeen     *prometheus.Desc
}

func NewUsersCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *UsersCollector {
	return &UsersCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		usersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "users_total"),
			"Total number of users.",
			nil, nil,
		),
		usersOnline: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "users_online"),
			"Number of users currently online.",
			nil, nil,
		),
		lastSeen: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "user", "last_seen_timestamp_seconds"),
			"Last seen timestamp of a user in seconds.",
			[]string{"user", "type"}, nil,
		),
	}
}

func (c *UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.usersTotal
	ch <- c.usersOnline
	ch <- c.lastSeen
	ch <- scrape.Duration
	ch <- scrape.Success
}

func (c *UsersCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	users, err := c.client.GetUsers(ctx)
	if err != nil {
		duration := time.Since(start).Seconds()
		c.logger.Error("users collector failed", "err", err)
		ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "users")
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "users")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.usersTotal, prometheus.GaugeValue, float64(len(users)))

	for _, user := range users {
		if user.LastSeen != nil {
			ch <- prometheus.MustNewConstMetric(c.lastSeen, prometheus.GaugeValue, float64(*user.LastSeen)/1000.0, user.Username, user.Type)
		}
	}

	onlineUsers, err := c.client.GetOnlineUsers(ctx)
	if err != nil {
		c.logger.Error("online users failed", "err", err)
		duration := time.Since(start).Seconds()
		ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "users")
		ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 0, "users")
		return
	}

	ch <- prometheus.MustNewConstMetric(c.usersOnline, prometheus.GaugeValue, float64(len(onlineUsers)))

	duration := time.Since(start).Seconds()
	ch <- prometheus.MustNewConstMetric(scrape.Duration, prometheus.GaugeValue, duration, "users")
	ch <- prometheus.MustNewConstMetric(scrape.Success, prometheus.GaugeValue, 1, "users")
}
