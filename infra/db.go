package infra

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetUpDB() *gorm.DB {
	env := os.Getenv("ENV")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port =%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	var db *gorm.DB
	var err error

	// production mode
	if env == "prod" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		log.Println("Setup postgres database")
	} else {
		// debug
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("Setup Sqlite memory database")
	}

	if err != nil {
		log.Fatal(err)
	}

	return db
}
