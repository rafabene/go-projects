// Package stresstest contém as funcionalidades principais para execução de testes de carga.
package stresstest

import "fmt"

// Config representa a configuração do teste de carga a ser executado.
// Contém todos os parâmetros necessários para definir como o teste será realizado.
type Config struct {
	URL         string // URL do serviço web que será testado
	Requests    int    // Número total de requests HTTP que serão enviados
	Concurrency int    // Número máximo de requests simultâneos (concorrência)
}

// Validate verifica se a configuração fornecida é válida para execução do teste.
// Retorna erro se algum parâmetro estiver incorreto ou inconsistente.
func (c *Config) Validate() error {
	// Verifica se a URL foi fornecida
	if c.URL == "" {
		return fmt.Errorf("URL é obrigatória")
	}
	
	// Verifica se o número de requests é positivo
	if c.Requests <= 0 {
		return fmt.Errorf("número de requests deve ser maior que 0")
	}
	
	// Verifica se o nível de concorrência é positivo
	if c.Concurrency <= 0 {
		return fmt.Errorf("concorrência deve ser maior que 0")
	}
	
	// Verifica se a concorrência não excede o total de requests
	// (não faz sentido ter mais workers que requests)
	if c.Concurrency > c.Requests {
		return fmt.Errorf("concorrência não pode ser maior que o número total de requests")
	}
	
	// Configuração válida
	return nil
}