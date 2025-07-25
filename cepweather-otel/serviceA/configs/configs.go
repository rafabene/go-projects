package configs

import (
	"fmt"
	"log"
	"os"
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

func GetServiceEndPoint() (string, error) {
	LoadConfig()
	endpoint := os.Getenv("SERVICEB_ENDPOINT")
	if endpoint == "" {
		return "", fmt.Errorf("SERVICEB_ENDPOINT não encontrado nas variáveis de ambiente")
	}
	return endpoint, nil
}
