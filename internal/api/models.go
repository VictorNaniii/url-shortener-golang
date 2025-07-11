package api

import "time"

// ShortenRequest defines payload for shorten URL endpoint
// swagger:model ShortenRequest
// Example: {"url":"https://example.com","custom_alias":"my-sale","expiration":"2025-12-31T23:59:59Z","max_clicks":100,"utm_source":"newsletter","utm_medium":"email","utm_campaign":"summer_sale"}
type ShortenRequest struct {
	URL         string  `json:"url" example:"https://example.com" binding:"required"`
	CustomAlias string  `json:"custom_alias,omitempty" example:"my-sale"`
	Expiration  string  `json:"expiration,omitempty" example:"2025-12-31T23:59:59Z"`
	MaxClicks   *uint64 `json:"max_clicks,omitempty" example:"100"`
	UTMSource   string  `json:"utm_source,omitempty" example:"newsletter"`
	UTMMedium   string  `json:"utm_medium,omitempty" example:"email"`
	UTMCampaign string  `json:"utm_campaign,omitempty" example:"summer_sale"`
}

// ShortenResponse defines response for shorten URL endpoint
// swagger:model ShortenResponse
// Example: {"short_url":"qIhf8TFq"}
type ShortenResponse struct {
	ShortURL string `json:"short_url" example:"qIhf8TFq"`
}

// StatsResponse defines response for stats endpoint
// swagger:model StatsResponse
// Example: {"short_url":"qIhf8TFq","original_url":"https://example.com","click_count":10,"last_clicked_at":"2025-07-11T22:00:00Z"}
type StatsResponse struct {
	ShortURL      string     `json:"short_url" example:"qIhf8TFq"`
	OriginalURL   string     `json:"original_url" example:"https://example.com"`
	ClickCount    uint64     `json:"click_count" example:"10"`
	LastClickedAt *time.Time `json:"last_clicked_at" example:"2025-07-11T22:00:00Z"`
}

// RegisterRequest defines payload for user registration
// swagger:model RegisterRequest
// Example: {"username":"ivancik","password":"secret"}
type RegisterRequest struct {
	Username string `json:"username" example:"ivancik" bindding:"required"`
	Password string `json:"password" example:"secret" binding:"required"`
}

// LoginRequest defines payload for user login
// swagger:model LoginRequest
// Example: {"username":"ivancik","password":"secret"}
type LoginRequest struct {
	Username string `json:"username" example:"ivancik" binding:"required"`
	Password string `json:"password" example:"secret" binding:"required"`
}

// LoginResponse defines response for login endpoint
// swagger:model LoginResponse
// Example: {"token":"<jwt-token>"}
type LoginResponse struct {
	Token string `json:"token" example:"<jwt-token>"`
}

// ErrorResponse defines error response
// swagger:model ErrorResponse
// Example: {"message":"error description"}
type ErrorResponse struct {
	Message string `json:"message" example:"error description"`
}

// SuccessResponse defines a simple success message
// swagger:model SuccessResponse
// Example: {"message":"Registration successful"}
type SuccessResponse struct {
	Message string `json:"message" example:"Registration successful"`
}
