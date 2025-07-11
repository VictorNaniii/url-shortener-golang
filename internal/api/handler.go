package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
	"time"
	"url-shortener/internal/domain"
	"url-shortener/internal/middleware"
	"url-shortener/internal/service"
)

type URLHandler struct {
	service domain.URLService
}

func NewURLHandler(service domain.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// ShortenURL godoc
// @Summary      Shorten a URL with marketing options
// @Description  Create a shortened link for the given URL with optional custom alias, expiration, click limit, and UTM parameters
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        request  body   ShortenRequest   true  "Shorten request payload"
// @Success      200      {object} ShortenResponse
// @Failure      400      {object} ErrorResponse
// @Failure      500      {object} ErrorResponse
// @Router       /shorten [post]
func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL         string  `json:"url"`
		CustomAlias string  `json:"custom_alias,omitempty"`
		Expiration  string  `json:"expiration,omitempty"`
		MaxClicks   *uint64 `json:"max_clicks,omitempty"`
		UTMSource   string  `json:"utm_source,omitempty"`
		UTMMedium   string  `json:"utm_medium,omitempty"`
		UTMCampaign string  `json:"utm_campaign,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// Validate custom alias format
	if req.CustomAlias != "" {
		for _, r := range req.CustomAlias {
			if !(r == '-' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
				http.Error(w, "custom_alias must be alphanumeric or dash", http.StatusBadRequest)
				return
			}
		}
		// Check uniqueness
		if urlObj, err := h.service.GetStats(req.CustomAlias); err == nil && urlObj != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "custom_alias already taken"})
			return
		}
	}
	// Parse expiration
	var expPtr *time.Time
	if req.Expiration != "" {
		exp, err := time.Parse(time.RFC3339, req.Expiration)
		if err != nil || exp.Before(time.Now()) {
			http.Error(w, "Invalid expiration (must be RFC3339 and in the future)", http.StatusBadRequest)
			return
		}
		expPtr = &exp
	}
	// Validate max clicks
	if req.MaxClicks != nil && *req.MaxClicks == 0 {
		http.Error(w, "max_clicks must be > 0", http.StatusBadRequest)
		return
	}
	userID, _ := middleware.UserIDFromContext(r.Context())
	shortURL, err := h.service.ShortenWithOptions(req.URL, userID, req.CustomAlias, expPtr, req.MaxClicks, req.UTMSource, req.UTMMedium, req.UTMCampaign)
	if err != nil {
		if err.Error() == "custom alias already in use" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "custom_alias already taken"})
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// RedirectURL godoc
// @Summary      Redirect to original URL
// @Description  Redirects from a shortened token to the original URL and enforces expiration and click limit
// @Tags         urls
// @Param        shortURL  path   string           true  "Short URL token or custom alias"
// @Success      302      {string} string        "redirect URL"
// @Failure      404      {object} ErrorResponse
// @Failure      410      {object} ErrorResponse
// @Router       /{shortURL} [get]
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	originalURL, err := h.service.Redirect(shortURL)
	if err != nil {
		if err.Error() == "link expired" || err.Error() == "click limit reached" {
			http.Error(w, err.Error(), http.StatusGone)
		} else {
			http.Error(w, "URL not found", http.StatusNotFound)
		}
		return
	}
	http.Redirect(w, r, originalURL, http.StatusFound)
}

// StatsURL godoc
// @Summary      Get click statistics
// @Description  Returns click count and last click date for a shortened URL or alias
// @Tags         urls
// @Produce      json
// @Param        shortURL  path   string        true  "Short URL token or alias"
// @Success      200      {object} StatsResponse
// @Failure      404      {object} ErrorResponse
// @Router       /stats/{shortURL} [get]
func (h *URLHandler) StatsURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "shortURL")
	// Remove any URL encoding that might have been applied
	if len(shortURL) > 8 {
		// If URL is longer than expected, it might be URL encoded
		// Extract just the token part
		if strings.Contains(shortURL, "%") {
			http.Error(w, "Invalid short URL format. Use only the token (e.g., 'qIhf8TFq')", http.StatusBadRequest)
			return
		}
	}

	urlObj, err := h.service.GetStats(shortURL)
	if err != nil {
		http.Error(w, "Statistics not found", http.StatusNotFound)
		return
	}
	res := struct {
		ShortURL      string     `json:"short_url"`
		OriginalURL   string     `json:"original_url"`
		ClickCount    uint64     `json:"click_count"`
		LastClickedAt *time.Time `json:"last_clicked_at"`
	}{
		ShortURL:      urlObj.ShortenedURL,
		OriginalURL:   urlObj.OriginalURL,
		ClickCount:    urlObj.ClickCount,
		LastClickedAt: urlObj.LastClickedAt,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// Add UserHandler for registration, login, and user URLs

type UserHandler struct {
	service    *service.UserService
	urlService domain.URLService
}

func NewUserHandler(service *service.UserService, urlService domain.URLService) *UserHandler {
	return &UserHandler{service: service, urlService: urlService}
}

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body   RegisterRequest  true  "Registration payload"
// @Success      201      {object} SuccessResponse
// @Failure      400      {object} ErrorResponse
// @Router       /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	_, err := h.service.Register(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Registration successful"))
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body   LoginRequest     true  "Login payload"
// @Success      200      {object} LoginResponse
// @Failure      400      {object} ErrorResponse
// @Failure      401      {object} ErrorResponse
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	res := struct {
		Token string `json:"token"`
	}{Token: token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}

// GetUserURLs godoc
// @Summary      List user URLs
// @Description  Returns all shortened URLs for the authenticated user
// @Tags         users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200      {array}  model.URL
// @Failure      401      {object} ErrorResponse
// @Router       /user/urls [get]
func (h *UserHandler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	urls, err := h.urlService.GetURLsByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve URLs", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
