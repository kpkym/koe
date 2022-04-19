package domain

type DlDomain struct {
	Code string `gorm:"primaryKey"`
	Data string
}

type TrackDomain struct {
	Code string `gorm:"primaryKey"`
	Data string
}

type FileSystemTreeDomain struct {
	Code string `gorm:"primaryKey"`
	Data string
}
