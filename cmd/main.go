package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"url-shortener/config"
	"url-shortener/internal/api"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // allow all headers
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // preflight cache duration
	}))
	r.Post("/shorten", urlHandler.ShortenURL)
	r.Get("/{shortURL}", urlHandler.RedirectURL)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
