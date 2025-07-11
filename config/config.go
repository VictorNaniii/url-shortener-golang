package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"url-shortener/internal/model"
)

func ConfigDB(localhost, dbUser, dbPassword, dbName, dbPort string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		localhost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to the database")
	}

	if err := db.AutoMigrate(&model.URL{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	return db, nil
}
