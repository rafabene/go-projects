package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/models"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/services"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
	tracer   trace.Tracer
)

func init() {
	var err error
	tracer, err = tracing.NewTracer()
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
}

func HandleCepWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePostCep(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostCep(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HandleCepWeather")
	defer span.End()
	input := &models.CepInput{}
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
	weatherData, err := services.CallServiceB(ctx, input.Cep)
	if err != nil {
		http.Error(w, "failed to call service B", http.StatusInternalServerError)
		log.Printf("Error calling service B: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err = json.Marshal(weatherData)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
		return
	}
	w.Write(b)
}
