// Package main contém o ponto de entrada da aplicação stress-test CLI.
// Esta ferramenta permite realizar testes de carga em serviços web com controle de concorrência.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/rafabene/go-projects/stress-test/pkg/stresstest"
)

// Variáveis globais para armazenar os valores dos parâmetros CLI
var (
	url         string // URL do serviço a ser testado
	requests    int    // Número total de requests a serem enviados
	concurrency int    // Número de requests simultâneos
)

// rootCmd define o comando raiz da aplicação CLI usando Cobra
var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "Uma ferramenta CLI para testes de carga em serviços web",
	Long: `Stress Test é uma ferramenta CLI desenvolvida em Go para realizar testes de carga
em serviços web. Permite especificar URL, número de requests e nível de concorrência.`,
	RunE: runStressTest,
}

// init configura os flags/parâmetros da linha de comando
func init() {
	// Configura os flags com versões curtas e longas
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 0, "Número total de requests (obrigatório)")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Número de chamadas simultâneas")
	
	// Marca flags como obrigatórios
	rootCmd.MarkFlagRequired("url")
	rootCmd.MarkFlagRequired("requests")
}

// runStressTest executa o teste de carga com os parâmetros fornecidos
func runStressTest(cmd *cobra.Command, args []string) error {
	// Cria a configuração com os valores dos flags
	config := stresstest.Config{
		URL:         url,
		Requests:    requests,
		Concurrency: concurrency,
	}
	
	// Valida a configuração antes de prosseguir
	if err := config.Validate(); err != nil {
		return fmt.Errorf("configuração inválida: %w", err)
	}

	// Exibe informações do teste que será executado
	fmt.Printf("Iniciando teste de carga...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Total de requests: %d\n", config.Requests)
	fmt.Printf("Concorrência: %d\n", config.Concurrency)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Executa o teste de carga e obtém o relatório
	report := stresstest.Run(config)
	
	// Exibe o relatório final
	stresstest.PrintReport(report)
	
	return nil
}

// main é o ponto de entrada da aplicação
func main() {
	// Executa o comando raiz e trata erros
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}