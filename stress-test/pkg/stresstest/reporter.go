package stresstest

import "fmt"

// PrintReport exibe o relatório final do teste de carga de forma formatada e amigável.
// Mostra métricas importantes como tempo total, throughput, códigos de status e taxa de erro.
// O relatório é exibido no terminal com emojis e formatação visual para facilitar leitura.
func PrintReport(report Report) {
	// Cabeçalho do relatório com separadores visuais
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📊 RELATÓRIO DO TESTE DE CARGA")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Métricas principais do teste
	fmt.Printf("⏱️  Tempo total gasto: %v\n", report.TotalTime)
	fmt.Printf("📨 Total de requests realizados: %d\n", report.TotalRequests)
	fmt.Printf("✅ Requests com status 200: %d\n", report.SuccessCount)
	
	// Mostra contagem de erros apenas se houver algum
	if report.ErrorCount > 0 {
		fmt.Printf("❌ Requests com erro: %d\n", report.ErrorCount)
	}
	
	// Seção de distribuição de códigos de status HTTP
	fmt.Println("\n📈 Distribuição de códigos de status:")
	for statusCode, count := range report.StatusCodes {
		// Calcula a porcentagem de cada código de status
		percentage := float64(count) / float64(report.TotalRequests) * 100
		fmt.Printf("   %d: %d requests (%.1f%%)\n", statusCode, count, percentage)
	}
	
	// Calcula e exibe throughput (requests por segundo)
	if report.TotalRequests > 0 {
		requestsPerSecond := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("\n🚀 Requests por segundo: %.2f req/s\n", requestsPerSecond)
	}
	
	// Rodapé do relatório
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}