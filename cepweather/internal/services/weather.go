package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/rafabene/go-projects/cepweather/internal/models"
)

func GetWeatherData(localidade string) (models.WeatherData, error) {
	apiKey, err := getApiKey()
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to get API key: %w", err)
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, url.QueryEscape(localidade))
	log.Printf("Fetching weather data from URL: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherData{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.WeatherData{}, fmt.Errorf("failed to fetch Weather for %s: %s", localidade, resp.Status)
	}
	var weatherData models.WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return models.WeatherData{}, err
	}
	return weatherData, nil
}

func getApiKey() (string, error) {
	// Pega o caminho absoluto do diretório atual do código
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)

	// Monta o caminho completo até /configs/.env
	envPath := filepath.Join(basePath, "../..", "configs", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		return "", fmt.Errorf("failed to load .env file: %w", err)
	}
	apiKey := os.Getenv("WEATHER_APIKEY")
	if apiKey == "" {
		return "", fmt.Errorf("WEATHER_APIKEY not found in .env file")
	}
	return apiKey, nil
}
