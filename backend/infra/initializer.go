package infra

import (
	"github.com/joho/godotenv"
	"log"
)

func Initialize() {
	err := godotenv.Load("./backend/.env")
	if err != nil {
		log.Fatal(err)
	}
}
