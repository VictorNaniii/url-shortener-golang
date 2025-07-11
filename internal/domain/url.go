package domain

import "url-shortener/internal/model"

type URLRepository interface {
	Save(url *model.URL) error
	FindByShortURL(shortURL string) (*model.URL, error)
}
