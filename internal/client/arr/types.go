package arr

type SystemStatus struct {
	Version        string `json:"version"`
	StartTime      string `json:"startTime"`
	Branch         string `json:"branch"`
	RuntimeName    string `json:"runtimeName"`
	RuntimeVersion string `json:"runtimeVersion"`
}

type HealthCheck struct {
	Source  string `json:"source"`
	Type   string `json:"type"`
	Message string `json:"message"`
	WikiUrl string `json:"wikiUrl"`
}

type QueueResponse struct {
	TotalRecords int           `json:"totalRecords"`
	Records      []QueueRecord `json:"records"`
}

type QueueRecord struct {
	Status                string `json:"status"`
	TrackedDownloadStatus string `json:"trackedDownloadStatus"`
	TrackedDownloadState  string `json:"trackedDownloadState"`
}

type RootFolder struct {
	Path       string `json:"path"`
	FreeSpace  int64  `json:"freeSpace"`
	TotalSpace int64  `json:"totalSpace"`
}
