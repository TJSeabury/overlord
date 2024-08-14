package types

import (
	"time"

	"gorm.io/gorm"
)

type ErrorDetails struct {
	Domain   string `gorm:"not null" json:"domain"`
	Error    string `gorm:"not null" json:"errortext"`
	URL      string `gorm:"not null" json:"url"`
	Filename string `gorm:"not null" json:"filename"`
	Line     int    `gorm:"not null" json:"line"`
	Column   int    `gorm:"not null" json:"column"`
	Datetime string `gorm:"not null" json:"datetime"`
	OS       string `gorm:"not null" json:"os"`
	Browser  string `gorm:"not null" json:"browser"`
}

type ErrorDetailsModel struct {
	ID        int            `gorm:"primaryKey" json:"id" tstype:"number|null"`
	CreatedAt time.Time      `json:"created_at" tstype:"string|null"`
	UpdatedAt time.Time      `json:"updated_at" tstype:"string|null"`
	DeletedAt gorm.DeletedAt `gorm:"index" tstype:"null|string"`
	ErrorDetails
	WebPropertyID int `gorm:"not null" json:"web_property_id" tstype:"number|null"`
}
