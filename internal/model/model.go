package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // store hashed password
	CreatedAt time.Time `json:"created_at"`
}

type ShortURModel struct {
	ShortURModel string
	OriginalURL  string
	CreatedAt    time.Time
	UserID       int // reference to User
}
