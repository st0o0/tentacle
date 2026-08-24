package jellyfin

type SystemInfo struct {
	Version            string `json:"Version"`
	OperatingSystem    string `json:"OperatingSystem"`
	SystemArchitecture string `json:"SystemArchitecture"`
	HasPendingRestart  bool   `json:"HasPendingRestart"`
}

type User struct {
	Id               string `json:"Id"`
	Name             string `json:"Name"`
	LastLoginDate    string `json:"LastLoginDate"`
	LastActivityDate string `json:"LastActivityDate"`
}

type Session struct {
	Id                 string           `json:"Id"`
	UserName           string           `json:"UserName"`
	DeviceName         string           `json:"DeviceName"`
	Client             string           `json:"Client"`
	ApplicationVersion string           `json:"ApplicationVersion"`
	RemoteEndPoint     string           `json:"RemoteEndPoint"`
	NowPlayingItem     *NowPlayingItem  `json:"NowPlayingItem"`
	PlayState          *PlayState       `json:"PlayState"`
	TranscodingInfo    *TranscodingInfo `json:"TranscodingInfo"`
}

type NowPlayingItem struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	RunTimeTicks int64  `json:"RunTimeTicks"`
}

type PlayState struct {
	PositionTicks int64  `json:"PositionTicks"`
	PlayMethod    string `json:"PlayMethod"`
	IsPaused      bool   `json:"IsPaused"`
}

type TranscodingInfo struct {
	Bitrate                  int      `json:"Bitrate"`
	IsVideoDirect            bool     `json:"IsVideoDirect"`
	IsAudioDirect            bool     `json:"IsAudioDirect"`
	HardwareAccelerationType string   `json:"HardwareAccelerationType"`
	VideoCodec               string   `json:"VideoCodec"`
	AudioCodec               string   `json:"AudioCodec"`
	Container                string   `json:"Container"`
	Width                    int      `json:"Width"`
	Height                   int      `json:"Height"`
	CompletionPercentage     float64  `json:"CompletionPercentage"`
	Framerate                float64  `json:"Framerate"`
	TranscodeReasons         []string `json:"TranscodeReasons"`
}

type VirtualFolder struct {
	Name           string `json:"Name"`
	CollectionType string `json:"CollectionType"`
	ItemId         string `json:"ItemId"`
}

type ItemsResponse struct {
	Items            []Item `json:"Items"`
	TotalRecordCount int    `json:"TotalRecordCount"`
}

type Item struct {
	Id           string        `json:"Id"`
	Name         string        `json:"Name"`
	Type         string        `json:"Type"`
	Size         int64         `json:"Size"`
	DateCreated  string        `json:"DateCreated"`
	MediaSources []MediaSource `json:"MediaSources"`
}

type MediaSource struct {
	Container    string        `json:"Container"`
	MediaStreams  []MediaStream `json:"MediaStreams"`
}

type MediaStream struct {
	Type   string `json:"Type"`
	Codec  string `json:"Codec"`
	Width  int    `json:"Width"`
	Height int    `json:"Height"`
}

type ItemCounts struct {
	MovieCount      int `json:"MovieCount"`
	SeriesCount     int `json:"SeriesCount"`
	EpisodeCount    int `json:"EpisodeCount"`
	ArtistCount     int `json:"ArtistCount"`
	AlbumCount      int `json:"AlbumCount"`
	SongCount       int `json:"SongCount"`
	BookCount       int `json:"BookCount"`
	MusicVideoCount int `json:"MusicVideoCount"`
	TrailerCount    int `json:"TrailerCount"`
	BoxSetCount     int `json:"BoxSetCount"`
	ProgramCount    int `json:"ProgramCount"`
	ItemCount       int `json:"ItemCount"`
}

type ScheduledTask struct {
	Id                        string               `json:"Id"`
	Name                      string               `json:"Name"`
	Category                  string               `json:"Category"`
	State                     string               `json:"State"`
	CurrentProgressPercentage float64              `json:"CurrentProgressPercentage"`
	LastExecutionResult       *TaskExecutionResult `json:"LastExecutionResult"`
}

type TaskExecutionResult struct {
	Status       string `json:"Status"`
	StartTimeUtc string `json:"StartTimeUtc"`
	EndTimeUtc   string `json:"EndTimeUtc"`
}

type ActivityLogResponse struct {
	Items            []ActivityLogEntry `json:"Items"`
	TotalRecordCount int                `json:"TotalRecordCount"`
}

type ActivityLogEntry struct {
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	Date     string `json:"Date"`
	Severity string `json:"Severity"`
	UserId   string `json:"UserId"`
}

type PluginInfo struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Version   string `json:"Version"`
	Status    string `json:"Status"`
	HasUpdate bool   `json:"HasUpdate"`
}

type DevicesResponse struct {
	Items            []DeviceInfo `json:"Items"`
	TotalRecordCount int          `json:"TotalRecordCount"`
}

type DeviceInfo struct {
	Id               string `json:"Id"`
	Name             string `json:"Name"`
	LastUserName     string `json:"LastUserName"`
	DateLastActivity string `json:"DateLastActivity"`
	AppName          string `json:"AppName"`
	AppVersion       string `json:"AppVersion"`
}

type PlaybackActivity struct {
	Date      string  `json:"Date"`
	UserId    string  `json:"UserId"`
	UserName  string  `json:"UserName"`
	PlayCount int     `json:"PlayCount"`
	WatchTime float64 `json:"WatchTime"`
}
