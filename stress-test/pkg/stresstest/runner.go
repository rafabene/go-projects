package stresstest

import (
	"net/http"
	"sync"
	"time"
)

// Run executa o teste de carga conforme a configuração fornecida.
// Cria goroutines para fazer requests HTTP de forma concorrente e coleta os resultados.
// Retorna um relatório consolidado com todas as métricas do teste.
func Run(config Config) Report {
	// Marca o tempo de início do teste para calcular duração total
	startTime := time.Now()
	
	// Canal para receber os resultados de cada request individual
	results := make(chan Result, config.Requests)
	
	// WaitGroup para aguardar conclusão de todas as goroutines
	var wg sync.WaitGroup
	
	// Semáforo para controlar o número de requests simultâneos
	// Funciona como um pool de workers limitado pela concorrência
	semaphore := make(chan struct{}, config.Concurrency)
	
	// Cria uma goroutine para cada request que será feito
	for i := 0; i < config.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			// Adquire um slot no semáforo (bloqueia se atingir limite)
			semaphore <- struct{}{}
			// Libera o slot ao terminar a função
			defer func() { <-semaphore }()
			
			// Executa o request HTTP e envia resultado para o canal
			result := makeRequest(config.URL)
			results <- result
		}()
	}
	
	// Goroutine para fechar o canal de resultados quando todos requests terminarem
	go func() {
		wg.Wait()     // Aguarda todas as goroutines de request terminarem
		close(results) // Fecha o canal para sinalizar fim da coleta
	}()
	
	// Inicializa o relatório que será preenchido com os dados coletados
	report := Report{
		StatusCodes: make(map[int]int), // Mapa para contar códigos de status
	}
	
	// Processa cada resultado recebido do canal
	for result := range results {
		report.TotalRequests++
		
		// Se houve erro de rede/timeout, conta como erro
		if result.Error != nil {
			report.ErrorCount++
		} else {
			// Conta o código de status retornado
			report.StatusCodes[result.StatusCode]++
			
			// Conta requests bem-sucedidos (status 200)
			if result.StatusCode == 200 {
				report.SuccessCount++
			}
		}
	}
	
	// Calcula o tempo total decorrido do teste
	report.TotalTime = time.Since(startTime)
	
	return report
}

// makeRequest executa uma única requisição HTTP GET para a URL especificada.
// Mede o tempo de resposta e captura erros ou códigos de status.
// Retorna um Result com as informações da requisição.
func makeRequest(url string) Result {
	// Marca o tempo de início da requisição individual
	start := time.Now()
	
	// Cria cliente HTTP com timeout de 30 segundos
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Executa a requisição GET
	resp, err := client.Get(url)
	
	// Calcula quanto tempo a requisição levou
	duration := time.Since(start)
	
	// Se houve erro (timeout, DNS, conexão, etc.), retorna resultado com erro
	if err != nil {
		return Result{
			Duration: duration,
			Error:    err,
		}
	}
	
	// Fecha o body da resposta para evitar vazamento de recursos
	defer resp.Body.Close()
	
	// Retorna resultado bem-sucedido com código de status
	return Result{
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Error:      nil,
	}
}