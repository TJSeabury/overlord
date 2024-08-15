package types

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ErrorDetails struct {
	Domain     string `gorm:"not null" json:"domain"`
	ErrorText  string `gorm:"not null" json:"errorText"`
	URL        string `gorm:"not null" json:"url"`
	Filename   string `gorm:"not null" json:"filename"`
	Line       int    `gorm:"not null" json:"line"`
	Column     int    `gorm:"not null" json:"column"`
	Datetime   string `gorm:"not null" json:"datetime"`
	UserAgent  string `gorm:"not null" json:"userAgent"`
	StackTrace string `gorm:"not null" json:"stackTrace"`
}

func (e *ErrorDetails) SanitizeDomain() error {
	e.Domain = strings.TrimSpace(e.Domain)
	e.Domain = strings.TrimSuffix(e.Domain, ".")
	e.Domain = strings.TrimPrefix(e.Domain, ".")

	// Test against a regex pattern to ensure only valid cahracters are present
	if matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9-.]{1,61}[a-zA-Z0-9]$`, e.Domain); !matched {
		return errors.New("validation error: Domain is invalid")
	}
	return nil
}

func (e *ErrorDetails) SanitizeURL() error {
	if matched, _ := regexp.MatchString(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`, e.URL); !matched {
		return errors.New("validation error: URL is invalid")
	}
	return nil
}

func (e *ErrorDetails) SanitizeFilename() error {
	if matched, _ := regexp.MatchString(`^[^<>:"\\/|?*\x00-\x1F]+[^ .]$`, e.Filename); !matched {
		return errors.New("validation error: Filename is invalid")
	}
	return nil
}

func (e *ErrorDetails) SanitizeDatetime() error {
	e.Datetime = strings.TrimSpace(e.Datetime)
	dt, err := time.Parse(time.RFC3339, e.Datetime)
	if err != nil {
		return errors.New("validation error: Datetime is invalid")
	}
	e.Datetime = dt.Format(time.RFC3339)
	return nil
}

func (e *ErrorDetails) SanitizeUserAgent() error {
	// Test if the user agent mathes a common pattern, otherwise reject it.
	if matched, _ := regexp.MatchString(`^Mozilla\/5\.0 \(Linux; U; Android (\d+\.)?(\d+\.)?(\*|\d+); [a-z]{2}-[a-z]{2}; (AFTA|AFTN|AFTS|AFTB|AFTT|AFTM|AFTKMST12|AFTRS) Build\/([A-Z0-9]+)\) AppleWebKit\/(\d+\.)?(\*|\d+) \(KHTML, like Gecko\) Version\/4\.0 Mobile Safari\/(\d+\.)?(\*|\d+)$`, e.UserAgent); !matched {
		return errors.New("validation error: UserAgent is invalid")
	}

	return nil
}

func (e *ErrorDetails) ValidateFields() error {
	if e.Domain == "" {
		return errors.New("validation error: Domain is required")
	}
	if e.ErrorText == "" {
		return errors.New("validation error: ErrorText is required")
	}
	if e.URL == "" {
		return errors.New("validation error: URL is required")
	}
	if e.Filename == "" {
		return errors.New("validation error: Filename is required")
	}
	if e.Line == 0 {
		return errors.New("validation error: Line is required")
	}
	if e.Column == 0 {
		return errors.New("validation error: Column is required")
	}
	if e.Datetime == "" {
		return errors.New("validation error: Datetime is required")
	}
	if e.UserAgent == "" {
		return errors.New("validation error: UserAgent is required")
	}

	if err := e.SanitizeDomain(); err != nil {
		return err
	}
	if err := e.SanitizeURL(); err != nil {
		return err
	}
	if err := e.SanitizeFilename(); err != nil {
		return err
	}
	if err := e.SanitizeDatetime(); err != nil {
		return err
	}
	// if err := e.SanitizeUserAgent(); err != nil {
	// 	return err
	// }

	return nil
}

type ErrorDetailsModel struct {
	ID        int            `gorm:"primaryKey" json:"id" tstype:"number|null"`
	CreatedAt time.Time      `json:"created_at" tstype:"string|null"`
	UpdatedAt time.Time      `json:"updated_at" tstype:"string|null"`
	DeletedAt gorm.DeletedAt `gorm:"index" tstype:"null|string"`
	ErrorDetails
	WebPropertyID int `gorm:"not null" json:"web_property_id" tstype:"number|null"`
}
