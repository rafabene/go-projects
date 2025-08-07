package config

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestGodotenv_CarregarArquivoEspecifico(t *testing.T) {
	// Limpar variáveis de ambiente
	os.Clearenv()
	defer os.Clearenv()
	
	// Criar arquivo de teste temporário
	envTestContent := `# Arquivo de teste para godotenv
PORTA_SERVIDOR=9000
LIMITE_IP_POR_SEGUNDO=15
TEMPO_BLOQUEIO_IP=240
LIMITE_TOKEN_POR_SEGUNDO=150
TEMPO_BLOQUEIO_TOKEN=360

# Tokens personalizados para teste
TOKEN_LIMITE_token_teste=75
TOKEN_LIMITE_admin=500
TOKEN_LIMITE_guest=5`
	
	err := os.WriteFile(".env.test", []byte(envTestContent), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar .env.test temporário: %v", err)
	}
	defer os.Remove(".env.test")
	
	// Carregar arquivo de teste específico
	err = godotenv.Load(".env.test")
	if err != nil {
		t.Fatalf("Erro ao carregar .env.test: %v", err)
	}
	
	// Criar configuração
	config, err := CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	// Verificar se os valores do arquivo .env.test foram carregados
	if config.PortaServidor != 9000 {
		t.Errorf("Porta do arquivo .env.test não foi carregada: esperado 9000, obtido %d", config.PortaServidor)
	}
	
	if config.LimiteIPPorSegundo != 15 {
		t.Errorf("Limite IP do arquivo .env.test não foi carregado: esperado 15, obtido %d", config.LimiteIPPorSegundo)
	}
	
	if config.LimiteTokenPorSegundo != 150 {
		t.Errorf("Limite token do arquivo .env.test não foi carregado: esperado 150, obtido %d", config.LimiteTokenPorSegundo)
	}
	
	if config.TempoBloqueioIP != 240*time.Second {
		t.Errorf("Tempo bloqueio IP do arquivo .env.test não foi carregado: esperado 240s, obtido %v", config.TempoBloqueioIP)
	}
	
	if config.TempoBloqueioToken != 360*time.Second {
		t.Errorf("Tempo bloqueio token do arquivo .env.test não foi carregado: esperado 360s, obtido %v", config.TempoBloqueioToken)
	}
	
	// Verificar tokens personalizados
	expectedTokens := map[string]int{
		"token_teste": 75,
		"admin":       500,
		"guest":       5,
	}
	
	for tokenName, expectedLimit := range expectedTokens {
		if limit, exists := config.TokensPersonalizados[tokenName]; !exists || limit != expectedLimit {
			t.Errorf("Token personalizado %s incorreto: esperado %d, obtido %d (exists: %v)", 
				tokenName, expectedLimit, limit, exists)
		}
	}
}

func TestGodotenv_PrioridadeVariaveisAmbiente(t *testing.T) {
	// Limpar variáveis de ambiente
	os.Clearenv()
	
	// Definir variável de ambiente que deve ter prioridade
	os.Setenv("PORTA_SERVIDOR", "7777")
	defer os.Clearenv()
	
	// Criar um arquivo .env temporário para o teste
	envContent := `PORTA_SERVIDOR=9000`
	err := os.WriteFile(".env", []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Erro ao criar .env temporário: %v", err)
	}
	defer os.Remove(".env")
	
	// Carregar arquivo .env
	err = godotenv.Load(".env")
	if err != nil {
		t.Fatalf("Erro ao carregar .env: %v", err)
	}
	
	// Verificar que a variável de ambiente tem prioridade sobre o arquivo
	porta := os.Getenv("PORTA_SERVIDOR")
	if porta != "7777" {
		t.Errorf("Variável de ambiente deveria ter prioridade: esperado '7777', obtido '%s'", porta)
	}
	
	// Configuração deve usar a variável de ambiente
	config, err := CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	if config.PortaServidor != 7777 {
		t.Errorf("Configuração deveria usar variável de ambiente: esperado 7777, obtido %d", config.PortaServidor)
	}
}

