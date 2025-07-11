package model

import (
	"gorm.io/gorm"
	"time"
)

// URL represents a shortened URL in the database.
type URL struct {
	gorm.Model
	ShortenedURL string    `gorm:"uniqueIndex;not null" json:"shortened_url"` // Shortened URL
	OriginalURL  string    `gorm:"not null" json:"original_url"`              // Original URL
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`          // Timestamp of creation
}

// TableName overrides the default table name for URL.
func (URL) TableName() string {
	return "short_urls"
}
