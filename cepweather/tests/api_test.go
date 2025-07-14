package tests

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabene/go-projects/cepweather/internal/handlers"
	"github.com/rafabene/go-projects/cepweather/internal/models"
)

func TestNullBody(t *testing.T) {
	r := httptest.NewRequest("POST", "/api/v1/weather", nil)
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}

}

func TestInvalidMethod(t *testing.T) {
	r := httptest.NewRequest("PUT", "/api/v1/weather", nil)
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestInvalidInput(t *testing.T) {
	jsonStr := `{"nome":"Rafael"` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestValidInputNoCeP(t *testing.T) {
	jsonStr := `{"nome":"Rafael"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepIncomplete(t *testing.T) {
	jsonStr := `{"cep":"75110"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepWrongSize(t *testing.T) {
	jsonStr := `{"cep":"751104301"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepWrongFormat(t *testing.T) {
	jsonStr := `{"cep":"75110-43"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestValidInput(t *testing.T) {
	jsonStr := `{"cep":"75110430"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	var weatherOutput models.WeatherOutput
	if err := json.NewDecoder(w.Body).Decode(&weatherOutput); err != nil {
		t.Errorf("Failed to decode response: %v", err)
		return
	}
	log.Printf("Weather Output: %+v", weatherOutput)
	if weatherOutput.TempC == 0 && weatherOutput.TempF == 0 && weatherOutput.TempK == 0 {
		t.Error("Expected non-zero temperature values, got zero values")
	}
}

func TestInvalidCep(t *testing.T) {
	jsonStr := `{"cep":"75123456"}` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/weather", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}
