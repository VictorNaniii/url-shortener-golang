package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/url"
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
	urlObj, err := s.repo.FindByShortURL(shortURL)
	if err != nil {
		// also try custom alias
		urlObj, err = s.repo.FindByCustomAlias(shortURL)
		if err != nil {
			return "", err
		}
	}
	// Check expiration
	if urlObj.Expiration != nil && time.Now().After(*urlObj.Expiration) {
		return "", errors.New("link expired")
	}
	// Check click limit
	if urlObj.MaxClicks != nil && urlObj.ClickCount >= *urlObj.MaxClicks {
		return "", errors.New("click limit reached")
	}
	// update click statistics
	urlObj.ClickCount++
	now := time.Now()
	urlObj.LastClickedAt = &now
	err = s.repo.Update(urlObj)
	if err != nil {
		return "", err
	}
	return urlObj.OriginalURL, nil
}

func (s *urlService) ShortenForUser(originalURL string, userID uint) (string, error) {
	// First check if this URL already exists for this user
	if userID > 0 {
		existingURLs, err := s.repo.GetURLsByUser(userID)
		if err == nil {
			for _, url := range existingURLs {
				if url.OriginalURL == originalURL {
					return url.ShortenedURL, nil // Return existing shortened URL
				}
			}
		}
	}

	shortURL := s.generateShortURL(originalURL)

	// Check if this short URL already exists, if so, generate a new one
	for {
		existing, err := s.repo.FindByShortURL(shortURL)
		if err != nil || existing == nil {
			break // Short URL is unique, we can use it
		}
		// If existing URL is the same, return it
		if existing.OriginalURL == originalURL {
			return existing.ShortenedURL, nil
		}
		// Generate a new short URL with timestamp to ensure uniqueness
		shortURL = s.generateShortURLWithTimestamp(originalURL)
	}

	url := &model.URL{
		OriginalURL:  originalURL,
		ShortenedURL: shortURL,
		CreatedAt:    time.Now(),
		UserID:       userID,
	}
	if err := s.repo.Save(url); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *urlService) ShortenWithOptions(originalURL string, userID uint, customAlias string, expiration *time.Time, maxClicks *uint64, utmSource, utmMedium, utmCampaign string) (string, error) {
	// Validate custom alias
	var shortURL string
	if customAlias != "" {
		// check format
		if !isValidAlias(customAlias) {
			return "", errors.New("invalid custom alias format")
		}
		// ensure uniqueness
		if _, err := s.repo.FindByCustomAlias(customAlias); err == nil {
			return "", errors.New("custom alias already in use")
		}
		shortURL = customAlias
	} else {
		shortURL = s.generateShortURL(originalURL)
	}
	// Append UTM params
	parsed, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	q := parsed.Query()
	if utmSource != "" && q.Get("utm_source") == "" {
		q.Set("utm_source", utmSource)
	}
	if utmMedium != "" && q.Get("utm_medium") == "" {
		q.Set("utm_medium", utmMedium)
	}
	if utmCampaign != "" && q.Get("utm_campaign") == "" {
		q.Set("utm_campaign", utmCampaign)
	}
	parsed.RawQuery = q.Encode()
	finalURL := parsed.String()
	// Prepare model
	url := &model.URL{
		OriginalURL:  finalURL,
		ShortenedURL: shortURL,
		CreatedAt:    time.Now(),
		UserID:       userID,
		CustomAlias:  customAlias,
		Expiration:   expiration,
		MaxClicks:    maxClicks,
		UTMSource:    utmSource,
		UTMMedium:    utmMedium,
		UTMCampaign:  utmCampaign,
	}
	// Save
	if err := s.repo.Save(url); err != nil {
		return "", err
	}
	return shortURL, nil
}

func (s *urlService) GetURLsByUser(userID uint) ([]model.URL, error) {
	return s.repo.GetURLsByUser(userID)
}

func (s *urlService) GetStats(shortURL string) (*model.URL, error) {
	return s.repo.FindByShortURL(shortURL)
}

func (s *urlService) generateShortURL(originalURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalURL))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

func (s *urlService) generateShortURLWithTimestamp(originalURL string) string {
	hasher := sha1.New()
	hasher.Write([]byte(originalURL + time.Now().String()))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:8]
}

// isValidAlias checks alias strings (alphanumeric, dash, underscore)
func isValidAlias(alias string) bool {
	for _, r := range alias {
		if !(r == '-' || r == '_' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return false
		}
	}
	return true
}
