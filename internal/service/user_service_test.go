package service_test

import (
	"testing"
	//"time"

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
	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)
	return db
}

func TestRegisterAndLogin(t *testing.T) {
	db := setupDB(t)
	repo := repository.NewUserRepository(db)
	us := service.NewUserService(repo)

	// Register a new user
	user, err := us.Register("testuser", "password123")
	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)

	// Duplicate registration should fail
	_, err = us.Register("testuser", "password123")
	assert.Error(t, err)

	// Successful login
	token, err := us.Login("testuser", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Invalid password
	_, err = us.Login("testuser", "wrongpass")
	assert.Error(t, err)

	// Non-existent user
	_, err = us.Login("nouser", "password")
	assert.Error(t, err)

	// Login token should expire in ~24h
	// Extract claims to verify exp (basic check)
	// (You can decode JWT if needed)
}
