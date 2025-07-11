package service

import (
	"crypto/sha1"
	"encoding/base64"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/model"
)

type urlService struct {
	repo domain.URLRepository
}

func NewURLService(repo domain.URLRepository) domain.URLService {
	return &urlService{repo: repo}
}

func (s *urlService) Shorten(originalURL string) (string, error) {
	shortURL := s.generateShortURL(originalURL)
	url := &model.URL{
		OriginalURL:  originalURL,
		ShortenedURL: shortURL,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.Save(url); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *urlService) Redirect(shortURL string) (string, error) {
	url, err := s.repo.FindByShortURL(shortURL)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}

func (s *urlService) generateShortURL(originalURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalURL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}
