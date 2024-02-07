package main

import (
	"gin-market/infra"
	"gin-market/models"
	"log"
)

// アプリケーションの起動と、DBマイグレーションを分ける

func main() {
	infra.Initialize()
	db := infra.SetUpDB()

	err := db.AutoMigrate(&models.Item{}, &models.User{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("マイグレーション完了")
}
