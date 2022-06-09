package domain

import "gorm.io/datatypes"

type WorkDomain struct {
	Code            string         `json:"code" gorm:"primaryKey"`
	Circle          string         `json:"circle"`
	DlCount         int64          `json:"dl_count"`
	Age             string         `json:"age"`
	Price           int64          `json:"price"`
	RateAverage2Dp  float64        `json:"rate_average_2dp"`
	RateCount       int64          `json:"rate_count"`
	Release         string         `json:"release"`
	ReviewCount     int64          `json:"review_count"`
	Title           string         `json:"title"`
	UserRating      float64        `json:"user_rating"`
	Vas             datatypes.JSON `json:"vas"`
	Tags            datatypes.JSON `json:"tags"`
	RateCountDetail datatypes.JSON `json:"rate_count_detail"`
	Rank            datatypes.JSON `json:"rank"`
}
