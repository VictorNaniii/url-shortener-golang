package repository

import (
	"gorm.io/gorm"
	"url-shortener/internal/domain"
	"url-shortener/internal/model"
)

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) domain.URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Save(url *model.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) FindByShortURL(shortURL string) (*model.URL, error) {
	var url model.URL
	if err := r.db.Where("shortened_url = ?", shortURL).First(&url).Error; err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) GetURLsByUser(userID uint) ([]model.URL, error) {
	var urls []model.URL
	if err := r.db.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		return nil, err
	}
	return urls, nil
}
