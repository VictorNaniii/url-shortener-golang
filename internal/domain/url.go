package domain

import "url-shortener/internal/model"

type URLRepository interface {
	Save(url *model.URL) error
	FindByShortURL(shortURL string) (*model.URL, error)
	GetURLsByUser(userID uint) ([]model.URL, error)
	Update(url *model.URL) error
}
