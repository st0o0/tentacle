package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type ScrapeDescs struct {
	Duration *prometheus.Desc
	Success  *prometheus.Desc
}

func NewScrapeDescs(namespace string) ScrapeDescs {
	return ScrapeDescs{
		Duration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
			"Duration of a collector scrape.",
			[]string{"collector"}, nil,
		),
		Success: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "success"),
			"Whether a collector scrape was successful.",
			[]string{"collector"}, nil,
		),
	}
}
