package jellyfin

import (
	"time"

	"github.com/st0o0/tentacle/internal/collector"
)

const namespace = "jellyfin"

var scrape = collector.NewScrapeDescs(namespace)

func parseJellyfinTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}
