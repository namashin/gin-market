package infra

import (
	"github.com/joho/godotenv"
)

func Initialize() error {
	return godotenv.Load("./backend/.env")
}
