// @title URL Shortener API
// @version 1.0
// @description Shorten URLs and view analytics
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"url-shortener/config"
	"url-shortener/docs"
	"url-shortener/internal/api"
	"url-shortener/internal/middleware"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Swagger info
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	localhost := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	db, err := config.ConfigDB(localhost, user, password, dbname, port)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo)
	urlHandler := api.NewURLHandler(urlService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService, urlService)

	r := chi.NewRouter()
	// Register middleware before routes
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // allow all headers
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // preflight cache duration
	}))
	// Swagger UI
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{shortURL}", urlHandler.RedirectURL)
	r.Get("/stats/{shortURL}", urlHandler.StatsURL)

	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)
	r.With(middleware.AuthMiddleware).Get("/user/urls", userHandler.GetUserURLs)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
