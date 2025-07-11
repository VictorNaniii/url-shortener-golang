package model

import (
	//"gorm.io/gorm"
	"time"
)

// swagger:model URL
type URL struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ShortenedURL  string     `gorm:"uniqueIndex;not null" json:"shortened_url"`          // Shortened URL
	OriginalURL   string     `gorm:"not null" json:"original_url"`                       // Original URL
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`                   // Timestamp of creation
	UserID        uint       `gorm:"index" json:"user_id"`                               // User association
	ClickCount    uint64     `gorm:"default:0" json:"click_count"`                       // Number of clicks
	LastClickedAt *time.Time `json:"last_clicked_at"`                                    // Timestamp of last click
	CustomAlias   string     `gorm:"uniqueIndex;size:255" json:"custom_alias,omitempty"` // Optional custom alias
	Expiration    *time.Time `json:"expiration,omitempty"`                               // Optional link expiration
	MaxClicks     *uint64    `json:"max_clicks,omitempty"`                               // Optional click limit
	UTMSource     string     `gorm:"size:255" json:"utm_source,omitempty"`               // Optional UTM source
	UTMMedium     string     `gorm:"size:255" json:"utm_medium,omitempty"`               // Optional UTM medium
	UTMCampaign   string     `gorm:"size:255" json:"utm_campaign,omitempty"`             // Optional UTM campaign
}

// TableName overrides the default table name for URL.
func (URL) TableName() string {
	return "short_urls"
}
