package jellyfin

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/st0o0/tentacle/internal/client/jellyfin"
)

type usersCollector struct {
	client  *jellyfin.Client
	timeout time.Duration
	logger  *slog.Logger

	total        *prometheus.Desc
	lastLogin    *prometheus.Desc
	lastActivity *prometheus.Desc
}

func newUsersCollector(client *jellyfin.Client, timeout time.Duration, logger *slog.Logger) *usersCollector {
	return &usersCollector{
		client:  client,
		timeout: timeout,
		logger:  logger,
		total: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "users", "total"),
			"Total number of Jellyfin users.",
			nil, nil,
		),
		lastLogin: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "user", "last_login_timestamp_seconds"),
			"Unix timestamp of the user's last login.",
			[]string{"user"}, nil,
		),
		lastActivity: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "user", "last_activity_timestamp_seconds"),
			"Unix timestamp of the user's last activity.",
			[]string{"user"}, nil,
		),
	}
}

func (c *usersCollector) Name() string { return "users" }

func (c *usersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.total
	ch <- c.lastLogin
	ch <- c.lastActivity
}

func (c *usersCollector) Update(ch chan<- prometheus.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	users, err := c.client.GetUsers(ctx)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(c.total, prometheus.GaugeValue, float64(len(users)))

	for _, u := range users {
		if t, ok := parseJellyfinTime(u.LastLoginDate); ok {
			ch <- prometheus.MustNewConstMetric(c.lastLogin, prometheus.GaugeValue, float64(t.Unix()), u.Name)
		}
		if t, ok := parseJellyfinTime(u.LastActivityDate); ok {
			ch <- prometheus.MustNewConstMetric(c.lastActivity, prometheus.GaugeValue, float64(t.Unix()), u.Name)
		}
	}
	return nil
}
