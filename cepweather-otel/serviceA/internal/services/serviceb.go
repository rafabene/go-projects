package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rafabene/go-projects/cepweather-otel/serviceA/configs"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/models"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer trace.Tracer
)

func init() {
	var err error
	tracer, err = tracing.NewTracer()
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
}

func CallServiceB(ctx context.Context, cep string) (*models.WeatherData, error) {
	ctx, span := tracer.Start(ctx, "CallServiceB")
	defer span.End()

	cepInput := models.CepInput{
		Cep: cep,
	}
	ep, err := configs.GetServiceEndPoint()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter endpoint do serviço: %w", err)
	}
	b, err := json.Marshal(cepInput)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar entrada do CEP: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para o serviço B: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar o serviço B: %w", err)
	}
	defer resp.Body.Close()
	output := &models.WeatherData{}
	json.NewDecoder(resp.Body).Decode(&output)
	return output, nil
}
