package dto

type SearchResponse struct {
	Pagination `json:"pagination"`
	Works      []Work `json:"works"`
}

type Pagination struct {
	CurrentPage int `json:"currentPage,omitempty"`
	PageSize    int `json:"pageSize,omitempty"`
	TotalCount  int `json:"totalCount,omitempty"`
}
type Work struct {
	Circle         string   `json:"circle"`
	CreateDate     string   `json:"create_date"`
	DlCount        int64    `json:"dl_count"`
	HasSubtitle    bool     `json:"has_subtitle"`
	Id             int      `json:"id"`
	Name           string   `json:"name"`
	Nsfw           bool     `json:"nsfw"`
	Price          int64    `json:"price"`
	RateAverage2Dp float64  `json:"rate_average_2dp"`
	RateCount      int      `json:"rate_count"`
	Release        string   `json:"release"`
	ReviewCount    int64    `json:"review_count"`
	Tags           []string `json:"tags"`
	Title          string   `json:"title"`
	UserRating     int      `json:"userRating"`
	Vas            []string `json:"vas"`
}
