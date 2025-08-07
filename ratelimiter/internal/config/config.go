// Package config gerencia o carregamento e validação de configurações da aplicação.
//
// Este pacote centraliza toda a lógica de configuração, carregando valores
// de arquivos .env e variáveis de ambiente com fallback para valores padrão.
// Utiliza a biblioteca godotenv para compatibilidade com formatos padrão.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config contém todas as configurações necessárias para o funcionamento da aplicação.
//
// A estrutura organiza as configurações em grupos lógicos:
// - Servidor: configurações de porta e rede
// - Rate Limiting: limites globais para IP e tokens
// - Tokens Personalizados: limites específicos por token de acesso
//
// Todas as configurações podem ser definidas via variáveis de ambiente,
// arquivo .env, ou usar valores padrão seguros para desenvolvimento.
type Config struct {
	// Configurações do servidor HTTP
	PortaServidor int // Porta onde o servidor irá escutar (padrão: 8080)
	
	// Configurações de rate limiting por endereço IP
	LimiteIPPorSegundo int           // Máximo de requisições por IP por segundo (padrão: 10)
	TempoBloqueioIP    time.Duration // Tempo de bloqueio quando IP excede limite (padrão: 5min)
	
	// Configurações de rate limiting por token de acesso
	LimiteTokenPorSegundo int           // Máximo de requisições por token por segundo (padrão: 100) 
	TempoBloqueioToken    time.Duration // Tempo de bloqueio quando token excede limite (padrão: 5min)
	
	// Limites personalizados por token específico
	// Chave: nome do token, Valor: limite personalizado em req/segundo
	TokensPersonalizados map[string]int
}

// CarregarConfig carrega e valida todas as configurações da aplicação.
//
// A função segue uma ordem de prioridade para obter as configurações:
//  1. Variáveis de ambiente do sistema (maior prioridade)
//  2. Arquivo .env (se existir)
//  3. Valores padrão (menor prioridade)
//
// O carregamento é tolerante a falhas - se o arquivo .env não existir
// ou houver valores inválidos, a aplicação continua com valores padrão.
//
// Retorna uma instância de Config totalmente inicializada e pronta para uso.
func CarregarConfig() (*Config, error) {
	config := &Config{
		TokensPersonalizados: make(map[string]int),
	}
	
	// Tenta carregar arquivo .env usando godotenv
	// Falha silenciosa para permitir execução sem arquivo .env
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Aviso: Não foi possível carregar .env: %v\n", err)
	}
	
	// Carrega configurações com fallback para valores padrão seguros
	config.PortaServidor = obterIntEnv("PORTA_SERVIDOR", 8080)
	config.LimiteIPPorSegundo = obterIntEnv("LIMITE_IP_POR_SEGUNDO", 10)
	config.TempoBloqueioIP = time.Duration(obterIntEnv("TEMPO_BLOQUEIO_IP", 300)) * time.Second
	config.LimiteTokenPorSegundo = obterIntEnv("LIMITE_TOKEN_POR_SEGUNDO", 100)
	config.TempoBloqueioToken = time.Duration(obterIntEnv("TEMPO_BLOQUEIO_TOKEN", 300)) * time.Second
	
	// Carrega configurações de tokens personalizados
	config.carregarTokensPersonalizados()
	
	return config, nil
}


// obterIntEnv obtém um valor inteiro de variável de ambiente com fallback seguro.
//
// A função tenta converter a variável de ambiente para inteiro e,
// em caso de falha (valor inválido ou ausente), retorna o valor padrão.
// Logs de aviso são exibidos para valores inválidos para facilitar debugging.
//
// Parâmetros:
//   - chave: nome da variável de ambiente
//   - valorPadrao: valor a ser usado se a variável não existir ou for inválida
//
// Retorna o valor convertido ou o padrão em caso de erro.
func obterIntEnv(chave string, valorPadrao int) int {
	valorStr := os.Getenv(chave)
	if valorStr == "" {
		return valorPadrao
	}
	
	valor, err := strconv.Atoi(valorStr)
	if err != nil {
		fmt.Printf("Aviso: Valor inválido para %s: %s. Usando padrão: %d\n", chave, valorStr, valorPadrao)
		return valorPadrao
	}
	
	return valor
}

// carregarTokensPersonalizados descobre e carrega limites específicos por token.
//
// A função percorre todas as variáveis de ambiente procurando por padrões
// do tipo "TOKEN_LIMITE_<nome>" e extrai o limite personalizado para cada token.
//
// Formato esperado:
//   TOKEN_LIMITE_abc123=200     -> token "abc123" com limite 200 req/s
//   TOKEN_LIMITE_vip_user=1000  -> token "vip_user" com limite 1000 req/s
//
// Valores inválidos são ignorados com log de aviso, permitindo que a
// aplicação continue funcionando com outros tokens válidos.
func (c *Config) carregarTokensPersonalizados() {
	// Percorre todas as variáveis de ambiente do sistema
	for _, env := range os.Environ() {
		partes := strings.SplitN(env, "=", 2)
		if len(partes) != 2 {
			continue
		}
		
		chave := partes[0]
		valor := partes[1]
		
		// Verifica se é uma configuração de token personalizado
		if strings.HasPrefix(chave, "TOKEN_LIMITE_") {
			// Extrai o nome do token removendo o prefixo
			nomeToken := strings.TrimPrefix(chave, "TOKEN_LIMITE_")
			
			// Tenta converter o limite para inteiro
			limite, err := strconv.Atoi(valor)
			if err != nil {
				fmt.Printf("Aviso: Limite inválido para token %s: %s\n", nomeToken, valor)
				continue
			}
			
			// Armazena o limite personalizado para este token
			c.TokensPersonalizados[nomeToken] = limite
		}
	}
}

// String retorna uma representação formatada das configurações carregadas.
//
// Útil para logging e debugging durante a inicialização da aplicação.
// Exibe todas as configurações de forma organizada e legível.
//
// A saída inclui:
// - Configurações do servidor (porta)
// - Limites globais de rate limiting
// - Tokens personalizados (se configurados)
func (c *Config) String() string {
	var sb strings.Builder
	
	sb.WriteString("=== Configuração do Rate Limiter ===\n")
	sb.WriteString(fmt.Sprintf("Porta do Servidor: %d\n", c.PortaServidor))
	sb.WriteString(fmt.Sprintf("Limite IP/segundo: %d\n", c.LimiteIPPorSegundo))
	sb.WriteString(fmt.Sprintf("Tempo Bloqueio IP: %v\n", c.TempoBloqueioIP))
	sb.WriteString(fmt.Sprintf("Limite Token/segundo: %d\n", c.LimiteTokenPorSegundo))
	sb.WriteString(fmt.Sprintf("Tempo Bloqueio Token: %v\n", c.TempoBloqueioToken))
	
	// Lista tokens personalizados se houver algum configurado
	if len(c.TokensPersonalizados) > 0 {
		sb.WriteString("Tokens Personalizados:\n")
		for token, limite := range c.TokensPersonalizados {
			sb.WriteString(fmt.Sprintf("  %s: %d req/s\n", token, limite))
		}
	}
	
	return sb.String()
}