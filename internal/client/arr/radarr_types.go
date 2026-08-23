package arr

type Movie struct {
	Title      string `json:"title"`
	Monitored  bool   `json:"monitored"`
	Status     string `json:"status"`
	HasFile    bool   `json:"hasFile"`
	SizeOnDisk int64  `json:"sizeOnDisk"`
	Year       int    `json:"year"`
}
