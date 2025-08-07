package config

import (
	"os"
	"testing"
	"time"
)

func TestCarregarConfig_ValoresPadrao(t *testing.T) {
	// Limpar variáveis de ambiente
	os.Clearenv()
	
	config, err := CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	// Verificar valores padrão
	if config.PortaServidor != 8080 {
		t.Errorf("Porta padrão incorreta: esperado 8080, obtido %d", config.PortaServidor)
	}
	
	if config.LimiteIPPorSegundo != 10 {
		t.Errorf("Limite IP padrão incorreto: esperado 10, obtido %d", config.LimiteIPPorSegundo)
	}
	
	if config.LimiteTokenPorSegundo != 100 {
		t.Errorf("Limite token padrão incorreto: esperado 100, obtido %d", config.LimiteTokenPorSegundo)
	}
	
	if config.TempoBloqueioIP != 300*time.Second {
		t.Errorf("Tempo bloqueio IP incorreto: esperado 300s, obtido %v", config.TempoBloqueioIP)
	}
	
	if config.TempoBloqueioToken != 300*time.Second {
		t.Errorf("Tempo bloqueio token incorreto: esperado 300s, obtido %v", config.TempoBloqueioToken)
	}
}

func TestCarregarConfig_VariaveisAmbiente(t *testing.T) {
	// Limpar variáveis de ambiente
	os.Clearenv()
	
	// Definir variáveis de ambiente de teste
	os.Setenv("PORTA_SERVIDOR", "9090")
	os.Setenv("LIMITE_IP_POR_SEGUNDO", "20")
	os.Setenv("TEMPO_BLOQUEIO_IP", "600")
	os.Setenv("LIMITE_TOKEN_POR_SEGUNDO", "200")
	os.Setenv("TEMPO_BLOQUEIO_TOKEN", "900")
	os.Setenv("TOKEN_LIMITE_test_token", "50")
	
	defer os.Clearenv()
	
	config, err := CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	// Verificar valores das variáveis de ambiente
	if config.PortaServidor != 9090 {
		t.Errorf("Porta incorreta: esperado 9090, obtido %d", config.PortaServidor)
	}
	
	if config.LimiteIPPorSegundo != 20 {
		t.Errorf("Limite IP incorreto: esperado 20, obtido %d", config.LimiteIPPorSegundo)
	}
	
	if config.LimiteTokenPorSegundo != 200 {
		t.Errorf("Limite token incorreto: esperado 200, obtido %d", config.LimiteTokenPorSegundo)
	}
	
	if config.TempoBloqueioIP != 600*time.Second {
		t.Errorf("Tempo bloqueio IP incorreto: esperado 600s, obtido %v", config.TempoBloqueioIP)
	}
	
	if config.TempoBloqueioToken != 900*time.Second {
		t.Errorf("Tempo bloqueio token incorreto: esperado 900s, obtido %v", config.TempoBloqueioToken)
	}
	
	// Verificar token personalizado
	if limite, exists := config.TokensPersonalizados["test_token"]; !exists || limite != 50 {
		t.Errorf("Token personalizado incorreto: esperado 50, obtido %d (exists: %v)", limite, exists)
	}
}

func TestCarregarConfig_ValoresInvalidos(t *testing.T) {
	// Limpar variáveis de ambiente
	os.Clearenv()
	
	// Definir variáveis inválidas
	os.Setenv("PORTA_SERVIDOR", "invalid")
	os.Setenv("LIMITE_IP_POR_SEGUNDO", "abc")
	os.Setenv("TOKEN_LIMITE_invalid_token", "xyz")
	
	defer os.Clearenv()
	
	config, err := CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	// Deve usar valores padrão para entradas inválidas
	if config.PortaServidor != 8080 {
		t.Errorf("Deveria usar valor padrão para porta: esperado 8080, obtido %d", config.PortaServidor)
	}
	
	if config.LimiteIPPorSegundo != 10 {
		t.Errorf("Deveria usar valor padrão para limite IP: esperado 10, obtido %d", config.LimiteIPPorSegundo)
	}
	
	// Token inválido não deve estar no mapa
	if _, exists := config.TokensPersonalizados["invalid_token"]; exists {
		t.Error("Token com valor inválido não deveria estar no mapa")
	}
}

func TestObterIntEnv(t *testing.T) {
	tests := []struct {
		nome        string
		chave       string
		valor       string
		valorPadrao int
		esperado    int
	}{
		{
			nome:        "Valor válido",
			chave:       "TEST_VALID",
			valor:       "42",
			valorPadrao: 10,
			esperado:    42,
		},
		{
			nome:        "Valor inválido",
			chave:       "TEST_INVALID",
			valor:       "abc",
			valorPadrao: 10,
			esperado:    10,
		},
		{
			nome:        "Variável inexistente",
			chave:       "TEST_NONEXISTENT",
			valor:       "",
			valorPadrao: 15,
			esperado:    15,
		},
	}
	
	for _, test := range tests {
		t.Run(test.nome, func(t *testing.T) {
			// Limpar variável
			os.Unsetenv(test.chave)
			
			// Definir se valor não for vazio
			if test.valor != "" {
				os.Setenv(test.chave, test.valor)
				defer os.Unsetenv(test.chave)
			}
			
			resultado := obterIntEnv(test.chave, test.valorPadrao)
			if resultado != test.esperado {
				t.Errorf("Esperado %d, obtido %d", test.esperado, resultado)
			}
		})
	}
}

func TestConfig_String(t *testing.T) {
	config := &Config{
		PortaServidor:         8080,
		LimiteIPPorSegundo:    10,
		TempoBloqueioIP:       5 * time.Minute,
		LimiteTokenPorSegundo: 100,
		TempoBloqueioToken:    5 * time.Minute,
		TokensPersonalizados: map[string]int{
			"vip":   200,
			"basic": 20,
		},
	}
	
	str := config.String()
	
	// Verificar se contém informações importantes
	expectedStrings := []string{
		"Rate Limiter",
		"8080",
		"10",
		"100",
		"5m0s",
		"vip",
		"basic",
		"200",
		"20",
	}
	
	for _, expected := range expectedStrings {
		if !contains(str, expected) {
			t.Errorf("String de configuração deveria conter '%s', mas não contém. String: %s", expected, str)
		}
	}
}

// Função auxiliar para verificar se uma string contém outra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}