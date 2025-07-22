package configs

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	// Pega o caminho absoluto do diretório atual do código
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	// Monta o caminho completo até /configs/.env
	envPath := filepath.Join(basePath, ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Error loading .env file. Will continue with env var: %v", err)
	}
}
