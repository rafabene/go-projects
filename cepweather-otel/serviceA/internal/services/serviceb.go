package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rafabene/go-projects/cepweather-otel/serviceA/configs"
	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/models"
)

func CallServiceB(cep string) (*models.WeatherData, error) {
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
	req, err := http.NewRequest("POST", ep, io.NopCloser(bytes.NewBuffer(b)))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição para o serviço B: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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
