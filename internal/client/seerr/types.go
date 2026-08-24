package seerr

type Status struct {
	Version   string `json:"version"`
	CommitTag string `json:"commitTag"`
}

type RequestCount struct {
	Total     int `json:"total"`
	Movie     int `json:"movie"`
	TV        int `json:"tv"`
	Pending   int `json:"pending"`
	Approved  int `json:"approved"`
	Available int `json:"available"`
	Declined  int `json:"declined"`
}

type UsersResponse struct {
	PageInfo PageInfo `json:"pageInfo"`
}

type PageInfo struct {
	Pages   int `json:"pages"`
	Results int `json:"results"`
}

type IssueCount struct {
	Total    int `json:"total"`
	Open     int `json:"open"`
	Resolved int `json:"resolved"`
}
