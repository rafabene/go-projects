package stresstest

import "time"

// Result representa o resultado de uma única requisição HTTP durante o teste de carga.
// Contém informações sobre o status, tempo de resposta e possíveis erros.
type Result struct {
	StatusCode int           // Código de status HTTP retornado (200, 404, 500, etc.)
	Duration   time.Duration // Tempo que a requisição levou para ser concluída
	Error      error         // Erro ocorrido durante a requisição, se houver
}

// Report contém o relatório consolidado de todo o teste de carga executado.
// Inclui métricas gerais, estatísticas de sucesso/erro e distribuição de status codes.
type Report struct {
	TotalTime      time.Duration // Tempo total gasto na execução de todo o teste
	TotalRequests  int           // Número total de requisições que foram executadas
	SuccessCount   int           // Quantidade de requisições que retornaram status 200
	StatusCodes    map[int]int   // Mapa com a distribuição de códigos de status (código -> quantidade)
	ErrorCount     int           // Número de requisições que falharam com erro de rede/timeout
}