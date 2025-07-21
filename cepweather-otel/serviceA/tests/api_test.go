package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internals/handlers"
)

func TestNullBody(t *testing.T) {
	r := httptest.NewRequest("POST", "/api/v1/cep", nil)
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}

}

func TestInvalidMethod(t *testing.T) {
	r := httptest.NewRequest("PUT", "/api/v1/cep", nil)
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestInvalidInput(t *testing.T) {
	jsonStr := `{"nome":"Rafael"` // Invalid JSON - missing closing brace
	r := httptest.NewRequest("POST", "/api/v1/cep", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestValidInputNoCeP(t *testing.T) {
	jsonStr := `{"nome":"Rafael"}` // Invalid JSON - missing cep property
	r := httptest.NewRequest("POST", "/api/v1/cep", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepIncomplete(t *testing.T) {
	jsonStr := `{"cep":"75110"}` // Invalid JSON - cep incomplete
	r := httptest.NewRequest("POST", "/api/v1/cep", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepWrongSize(t *testing.T) {
	jsonStr := `{"cep":"751104301"}` // Invalid JSON - extra digits in cep
	r := httptest.NewRequest("POST", "/api/v1/cep", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}

func TestCepWrongFormat(t *testing.T) {
	jsonStr := `{"cep":"75110-43"}` // Invalid JSON - invalid format
	r := httptest.NewRequest("POST", "/api/v1/cep", strings.NewReader(jsonStr))
	w := httptest.NewRecorder()

	handlers.HandleCepWeather(w, r)
	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
	}
}
