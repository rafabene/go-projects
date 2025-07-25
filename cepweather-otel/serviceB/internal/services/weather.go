package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/models"
	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/tracing"

	"go.opentelemetry.io/otel/trace"
)

var (
	tracerWeather trace.Tracer
)

func init() {
	var err error
	tracerWeather, err = tracing.NewTracer()
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
}

func GetWeatherData(ctx context.Context, localidade string) (models.WeatherData, error) {
	ctx, span := tracerWeather.Start(ctx, "GetWeather")
	defer span.End()

	apiKey, err := getApiKey()
	if err != nil {
		return models.WeatherData{}, fmt.Errorf("failed to get API key: %w", err)
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, url.QueryEscape(localidade))
	urlEscaped := strings.Replace(url, apiKey, "******", 1) // Mask the API key in logs
	log.Printf("Fetching weather data from URL: %s", urlEscaped)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.WeatherData{}, err
	}
	resp, err := http.DefaultClient.Do(req)
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
		log.Printf("Error loading .env file. Will continue with env var: %v", err)

	}
	apiKey := os.Getenv("WEATHER_APIKEY")
	if apiKey == "" {
		return "", fmt.Errorf("WEATHER_APIKEY not found in environment variables")
	}
	return apiKey, nil
}
