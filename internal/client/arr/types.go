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

type Backup struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
	Time string `json:"time"`
}

type Update struct {
	Version   string `json:"version"`
	Installed bool   `json:"installed"`
	Latest    bool   `json:"latest"`
}

type BlocklistResponse struct {
	TotalRecords int `json:"totalRecords"`
}

type DownloadClient struct {
	Name     string `json:"name"`
	Protocol string `json:"protocol"`
	Priority int    `json:"priority"`
	Enable   bool   `json:"enable"`
}
