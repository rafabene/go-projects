package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/models"
)

func GetCepData(cep string) (models.CepOutput, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", url.QueryEscape(cep))
	log.Printf("Fetching CEP data from URL: %s", url)
	resp, err := http.Get(url)
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
