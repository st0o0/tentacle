package arr

type Series struct {
	Title             string `json:"title"`
	Monitored         bool   `json:"monitored"`
	Status            string `json:"status"`
	SeasonCount       int    `json:"seasonCount"`
	EpisodeCount      int    `json:"episodeCount"`
	EpisodeFileCount  int    `json:"episodeFileCount"`
	TotalEpisodeCount int    `json:"totalEpisodeCount"`
	SizeOnDisk        int64  `json:"sizeOnDisk"`
}

type WantedResponse struct {
	TotalRecords int `json:"totalRecords"`
}

type CalendarEntry struct {
	Title string `json:"title"`
}
