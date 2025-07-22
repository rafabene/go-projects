package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/models"
	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/tracing"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracerCep trace.Tracer
)

func init() {
	var err error
	tracerCep, err = tracing.NewTracer()
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
	}
}

func GetCepData(ctx context.Context, cep string) (models.CepOutput, error) {
	ctx, span := tracerCep.Start(ctx, "GetCEP")
	defer span.End()
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", url.QueryEscape(cep))
	log.Printf("Fetching CEP data from URL: %s", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.CepOutput{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.CepOutput{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.CepOutput{}, fmt.Errorf("failed to fetch data for CEP %s: %s", cep, resp.Status)
	}
	var cepOutput models.CepOutput
	if err := json.NewDecoder(resp.Body).Decode(&cepOutput); err != nil {
		return models.CepOutput{}, err
	}
	if cepOutput.Cep == "" {
		return models.CepOutput{}, fmt.Errorf("invalid CEP: %s", cep)
	}
	return cepOutput, nil
}
