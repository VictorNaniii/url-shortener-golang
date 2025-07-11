package domain

import "url-shortener/internal/model"

type URLService interface {
	Shorten(originalURL string) (string, error)
	ShortenForUser(originalURL string, userID uint) (string, error)
	Redirect(shortURL string) (string, error)
	GetURLsByUser(userID uint) ([]model.URL, error)
	GetStats(shortURL string) (*model.URL, error)
}
