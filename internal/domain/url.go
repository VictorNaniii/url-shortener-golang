package domain

import (
	"time"
	"url-shortener/internal/model"
)

type URLRepository interface {
	Save(url *model.URL) error
	FindByShortURL(shortURL string) (*model.URL, error)
	FindByCustomAlias(alias string) (*model.URL, error)
	GetURLsByUser(userID uint) ([]model.URL, error)
	Update(url *model.URL) error
}

// URLService interface
type URLService interface {
	Shorten(originalURL string) (string, error)
	ShortenForUser(originalURL string, userID uint) (string, error)
	Redirect(shortURL string) (string, error)
	GetURLsByUser(userID uint) ([]model.URL, error)
	GetStats(shortURL string) (*model.URL, error)
	ShortenWithOptions(originalURL string, userID uint, customAlias string, expiration *time.Time, maxClicks *uint64, utmSource, utmMedium, utmCampaign string) (string, error)
}
