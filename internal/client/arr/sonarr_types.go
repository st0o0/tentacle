package arr

type SeriesStatistics struct {
	SeasonCount       int   `json:"seasonCount"`
	EpisodeCount      int   `json:"episodeCount"`
	EpisodeFileCount  int   `json:"episodeFileCount"`
	TotalEpisodeCount int   `json:"totalEpisodeCount"`
	SizeOnDisk        int64 `json:"sizeOnDisk"`
}

type Series struct {
	Title             string           `json:"title"`
	Monitored         bool             `json:"monitored"`
	Status            string           `json:"status"`
	SeasonCount       int              `json:"seasonCount"`
	EpisodeCount      int              `json:"episodeCount"`
	EpisodeFileCount  int              `json:"episodeFileCount"`
	TotalEpisodeCount int              `json:"totalEpisodeCount"`
	SizeOnDisk        int64            `json:"sizeOnDisk"`
	Statistics        SeriesStatistics `json:"statistics"`
}

type WantedResponse struct {
	TotalRecords int `json:"totalRecords"`
}

type CalendarEntry struct {
	Title string `json:"title"`
}
