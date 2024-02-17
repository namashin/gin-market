package infra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

	// production for postgres
	if env == "prod" {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		log.Println("Setup postgres database")
	} else {
		// debug for sqlite memory
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		log.Println("Setup Sqlite memory database")
	}

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func ConnectToRedis(ctx *gin.Context) (*redis.Client, error) {
	// Redisとの接続を確立し、クライアントを返す
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redisのアドレス
		Password: "",               // Redisのパスワード（必要に応じて設定）
		DB:       0,                // 使用するデータベース
	})

	// 接続をテストする
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
