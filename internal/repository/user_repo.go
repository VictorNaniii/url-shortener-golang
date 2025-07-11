package repository

import (
	"errors"
	"gorm.io/gorm"
	"url-shortener/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user model.User) (model.User, error) {
	var existing model.User
	if err := r.db.Where("username = ?", user.Username).First(&existing).Error; err == nil {
		return model.User{}, errors.New("username already exists")
	}
	if err := r.db.Create(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(id int) (model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}
