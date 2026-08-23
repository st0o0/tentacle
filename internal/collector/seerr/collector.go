package seerr

import "github.com/st0o0/tentacle/internal/collector"

const namespace = "seerr"

var scrape = collector.NewScrapeDescs(namespace)
