package service_test

import (
	"testing"
	"time"

	"url-shortener/internal/model"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&model.User{}, &model.URL{})
	assert.NoError(t, err)
	return db
}

func TestShortenURLAndRedirect(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)

	// Test shorten
	orig := "https://example.com"
	token, err := svc.ShortenForUser(orig, 0)
	assert.NoError(t, err)
	assert.Len(t, token, 8)

	// Verify stored
	urlObj, err := repo.FindByShortURL(token)
	assert.NoError(t, err)
	assert.Equal(t, orig, urlObj.OriginalURL)

	// Test redirect increments click count
	resURL, err := svc.Redirect(token)
	assert.NoError(t, err)
	assert.Equal(t, orig, resURL)

	stats, err := svc.GetStats(token)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), stats.ClickCount)
	assert.NotNil(t, stats.LastClickedAt)
	assert.WithinDuration(t, time.Now(), *stats.LastClickedAt, time.Minute)
}

func TestDuplicateShorten(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)

	orig := "https://duplicate.com"
	t1, err := svc.ShortenForUser(orig, 0)
	assert.NoError(t, err)
	t2, err := svc.ShortenForUser(orig, 0)
	assert.NoError(t, err)
	assert.Equal(t, t1, t2)

	// Only one record in DB
	var urls []model.URL
	db.Find(&urls)
	assert.Len(t, urls, 1)
}

func TestShortenForUser(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewURLRepository(db)
	svc := service.NewURLService(repo)

	userID := uint(42)
	orig1 := "https://user.com/page1"
	t1, err := svc.ShortenForUser(orig1, userID)
	assert.NoError(t, err)

	orig2 := "https://user.com/page2"
	t2, err := svc.ShortenForUser(orig2, userID)
	assert.NoError(t, err)
	assert.NotEqual(t, t1, t2)

	urls, err := svc.GetURLsByUser(userID)
	assert.NoError(t, err)
	assert.Len(t, urls, 2)
}
