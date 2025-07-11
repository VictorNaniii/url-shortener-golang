package model

import (
	"gorm.io/gorm"
	"time"
)

// URL represents a shortened URL in the database.
type URL struct {
	gorm.Model
	ShortenedURL  string     `gorm:"uniqueIndex;not null" json:"shortened_url"` // Shortened URL
	OriginalURL   string     `gorm:"not null" json:"original_url"`              // Original URL
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`          // Timestamp of creation
	UserID        uint       `gorm:"index" json:"user_id"`                      // User association
	ClickCount    uint64     `gorm:"default:0" json:"click_count"`              // Number of clicks
	LastClickedAt *time.Time `json:"last_clicked_at"`                           // Timestamp of last click
}

// TableName overrides the default table name for URL.
func (URL) TableName() string {
	return "short_urls"
}
