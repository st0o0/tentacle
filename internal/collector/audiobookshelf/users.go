package audiobookshelf

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/audiobookshelf"
)

type usersCollector struct {
	client  *audiobookshelf.Client
	timeout time.Duration
	logger  *slog.Logger

	usersTotal    *prometheus.Desc
	usersActive   *prometheus.Desc
	usersOnline   *prometheus.Desc
	lastSeen      *prometheus.Desc
	listeningTime *prometheus.Desc
}

func newUsersCollector(client *audiobookshelf.Client, timeout time.Duration, logger *slog.Logger) *usersCollector {
	return &usersCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		usersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "users_total"),
			"Total number of users.",
			nil, nil,
		),
		usersActive: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "users_active_total"),
			"Number of active users.",
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
		listeningTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "user", "listening_time_seconds"),
			"Total listening time per user in seconds.",
			[]string{"user"}, nil,
		),
	}
}

func (c *usersCollector) Name() string {
	return "users"
}

func (c *usersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.usersTotal
	ch <- c.usersActive
	ch <- c.usersOnline
	ch <- c.lastSeen
	ch <- c.listeningTime
}

func (c *usersCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	users, err := c.client.GetUsers(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.usersTotal, prometheus.GaugeValue, float64(len(users)))

	activeCount := 0
	for _, user := range users {
		if user.IsActive {
			activeCount++
		}
		if user.LastSeen != nil {
			ch <- prometheus.MustNewConstMetric(c.lastSeen, prometheus.GaugeValue, float64(*user.LastSeen)/1000.0, user.Username, user.Type)
		}
		stats, err := c.client.GetUserListeningStats(ctx, user.Id)
		if err != nil {
			c.logger.Warn("user listening stats failed", "user", user.Username, "err", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.listeningTime, prometheus.GaugeValue, stats.TotalTime, user.Username)
	}
	ch <- prometheus.MustNewConstMetric(c.usersActive, prometheus.GaugeValue, float64(activeCount))

	onlineUsers, err := c.client.GetOnlineUsers(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.usersOnline, prometheus.GaugeValue, float64(len(onlineUsers)))

	return nil
}
