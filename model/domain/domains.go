package domain

import "gorm.io/datatypes"

type WorkDomain struct {
	Circle         string         `json:"circle" gorm:"primaryKey"`
	CreateDate     string         `json:"create_date"`
	DlCount        int64          `json:"dl_count"`
	HasSubtitle    bool           `json:"has_subtitle"`
	Code           int            `json:"code"`
	Nsfw           bool           `json:"nsfw"`
	Price          int64          `json:"price"`
	RateAverage2Dp float64        `json:"rate_average_2dp"`
	RateCount      int64          `json:"rate_count"`
	Release        string         `json:"release"`
	ReviewCount    int64          `json:"review_count"`
	Tags           datatypes.JSON `json:"tags"`
	Title          string         `json:"title"`
	UserRating     int            `json:"userRating"`
	Vas            datatypes.JSON `json:"vas"`
}
