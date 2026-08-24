package arr

type Indexer struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Enable   bool   `json:"enable"`
	Protocol string `json:"protocol"`
	Priority int    `json:"priority"`
}

type IndexerStatsResponse struct {
	Indexers []IndexerStats `json:"indexers"`
}

type IndexerStats struct {
	IndexerName          string `json:"indexerName"`
	NumberOfQueries      int    `json:"numberOfQueries"`
	NumberOfGrabs        int    `json:"numberOfGrabs"`
	NumberOfFailedQueries int   `json:"numberOfFailedQueries"`
	NumberOfFailedGrabs  int    `json:"numberOfFailedGrabs"`
	AverageResponseTime  int    `json:"averageResponseTime"`
}

type Application struct {
	Name           string `json:"name"`
	SyncLevel      string `json:"syncLevel"`
	Implementation string `json:"implementation"`
}

type IndexerStatus struct {
	IndexerId    int    `json:"indexerId"`
	DisabledTill string `json:"disabledTill"`
}
