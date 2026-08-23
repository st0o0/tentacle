package audiobookshelf

type ServerStatus struct {
	IsInitialized bool   `json:"isInitialized"`
	Language      string `json:"language"`
	ConfigPath    string `json:"ConfigPath"`
	MetadataPath  string `json:"MetadataPath"`
}

type AuthResponse struct {
	User           User           `json:"user"`
	ServerSettings ServerSettings `json:"serverSettings"`
}

type ServerSettings struct {
	ScannerMaxThreads int  `json:"scannerMaxThreads"`
	BackupSchedule    bool `json:"backupSchedule"`
}

type Library struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	MediaType string `json:"mediaType"`
	Provider  string `json:"provider"`
	CreatedAt int64  `json:"createdAt"`
	LastUpdate int64 `json:"lastUpdate"`
}

type LibraryStats struct {
	TotalItems       int            `json:"totalItems"`
	TotalSize        int64          `json:"totalSize"`
	TotalDuration    float64        `json:"totalDuration"`
	NumAudioTracks   int            `json:"numAudioTracks"`
	TotalAuthors     int            `json:"totalAuthors"`
	TotalGenres      int            `json:"totalGenres"`
	NumMissing       int            `json:"numMissing"`
	NumInvalid       int            `json:"numInvalid"`
	LargestItems     []ItemMinified `json:"largestItems"`
	LongestItems     []ItemMinified `json:"longestItems"`
	AuthorsWithCount []AuthorCount  `json:"authorsWithCount"`
	GenresWithCount  []GenreCount   `json:"genresWithCount"`
}

type ItemMinified struct {
	Id       string  `json:"id"`
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	Size     int64   `json:"size"`
}

type AuthorCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type GenreCount struct {
	Genre string `json:"genre"`
	Count int    `json:"count"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Type     string `json:"type"`
	IsActive bool   `json:"isActive"`
	LastSeen *int64 `json:"lastSeen"`
	CreatedAt int64 `json:"createdAt"`
}

type ListeningStats struct {
	TotalTime      float64            `json:"totalTime"`
	Today          float64            `json:"today"`
	Days           map[string]float64 `json:"days"`
	DayOfWeek      map[string]float64 `json:"dayOfWeek"`
	RecentSessions []ListeningSession `json:"recentSessions"`
}

type ListeningSession struct {
	Id            string  `json:"id"`
	UserId        string  `json:"userId"`
	LibraryId     string  `json:"libraryId"`
	MediaType     string  `json:"mediaType"`
	Duration      float64 `json:"duration"`
	PlayMethod    int     `json:"playMethod"`
	StartedAt     int64   `json:"startedAt"`
	UpdatedAt     int64   `json:"updatedAt"`
	DisplayTitle  string  `json:"displayTitle"`
	DisplayAuthor string  `json:"displayAuthor"`
}

type SessionsResponse struct {
	Total        int                `json:"total"`
	NumPages     int                `json:"numPages"`
	Page         int                `json:"page"`
	ItemsPerPage int                `json:"itemsPerPage"`
	Sessions     []ListeningSession `json:"sessions"`
}

type OnlineUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Backup struct {
	Id            string `json:"id"`
	Path          string `json:"path"`
	Datestamp     string `json:"datePretty"`
	CreatedAt     int64  `json:"createdAt"`
	ServerVersion string `json:"serverVersion"`
}

type BackupsResponse struct {
	Backups []Backup `json:"backups"`
}

type LibraryItemsResponse struct {
	Results   []LibraryItem `json:"results"`
	Total     int           `json:"total"`
	Limit     int           `json:"limit"`
	Page      int           `json:"page"`
	SortBy    string        `json:"sortBy"`
	SortDesc  bool          `json:"sortDesc"`
	FilterBy  string        `json:"filterBy"`
	MediaType string        `json:"mediaType"`
	Minified  bool          `json:"minified"`
	Collapsed bool          `json:"collapsified"`
}

type LibraryItem struct {
	Id        string `json:"id"`
	MediaType string `json:"mediaType"`
}
