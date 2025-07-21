package handlers

import (
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/models"
	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/services"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func HandleCepWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostCepWeather(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func roundToTwoDecimals(f float64) float64 {
	return math.Round(f*100) / 100
}

func handlePostCepWeather(w http.ResponseWriter, r *http.Request) {
	input := &models.WeatherInput{}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	if err := json.Unmarshal(b, input); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(input); err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		log.Printf("Validation error: %v", err)
		return
	}
	cepData, err := services.GetCepData(input.Cep)
	if err != nil {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	weatherData, err := services.GetWeatherData(cepData.Localidade)
	if err != nil {
		http.Error(w, "can not find weather data", http.StatusInternalServerError)
		log.Printf("Error fetching weather data: %v", err)
		return
	}

	weatherOutput := models.WeatherOutput{
		TempC: weatherData.Current.TempC,
		TempF: weatherData.Current.TempF,
		TempK: roundToTwoDecimals(weatherData.Current.TempC + 273.15),
	}
	w.Header().Set("Content-Type", "application/json")
	b, err = json.Marshal(weatherOutput)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
	w.Write(b)
}
