package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/rafabene/go-projects/ratelimiter/internal/config"
	"github.com/rafabene/go-projects/ratelimiter/internal/middleware"
)

func TestIntegracao_RateLimiterCompleto(t *testing.T) {
	// Configurar servidor de teste
	servidor := criarServidorTeste()
	defer servidor.Close()
	
	// Teste de limitação por IP
	t.Run("LimitacaoIP", func(t *testing.T) {
		testeLimitacaoIP(t, servidor.URL)
	})
	
	// Teste de limitação por token
	t.Run("LimitacaoToken", func(t *testing.T) {
		testeLimitacaoToken(t, servidor.URL)
	})
	
	// Teste de concorrência
	t.Run("ConcorrenciaIP", func(t *testing.T) {
		testeConcorrenciaIP(t, servidor.URL)
	})
	
	// Teste de diferentes IPs
	t.Run("DiferentesIPs", func(t *testing.T) {
		testeDiferentesIPs(t, servidor.URL)
	})
}

func criarServidorTeste() *httptest.Server {
	
	// Configurar rate limiter para testes
	configRL := &middleware.ConfigRateLimiter{
		LimiteIPPorSegundo:    3,
		TempoBloqueioIP:       2 * time.Second,
		LimiteTokenPorSegundo: 5,
		TempoBloqueioToken:    2 * time.Second,
		TokensPersonalizados: map[string]int{
			"token_vip": 10,
			"token_basic": 2,
		},
	}
	
	rateLimiter := middleware.NovoRateLimiter(configRL)
	
	// Handler simples de teste
	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resposta := map[string]interface{}{
			"status": "sucesso",
			"timestamp": time.Now(),
			"path": r.URL.Path,
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resposta)
	}))
	
	return httptest.NewServer(handler)
}

func testeLimitacaoIP(t *testing.T, baseURL string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	// Fazer 3 requisições (limite)
	for i := 0; i < 3; i++ {
		resp, err := client.Get(baseURL + "/teste")
		if err != nil {
			t.Fatalf("Erro na requisição %d: %v", i+1, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Requisição %d deveria ser permitida, status: %d", i+1, resp.StatusCode)
		}
	}
	
	// Quarta requisição deve ser bloqueada
	resp, err := client.Get(baseURL + "/teste")
	if err != nil {
		t.Fatalf("Erro na quarta requisição: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Quarta requisição deveria ser bloqueada, status: %d", resp.StatusCode)
	}
	
	// Verificar header Retry-After
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		t.Error("Header Retry-After deveria estar presente")
	}
}

func testeLimitacaoToken(t *testing.T, baseURL string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	// Teste com token básico (limite 2)
	for i := 0; i < 2; i++ {
		req, _ := http.NewRequest("GET", baseURL+"/teste", nil)
		req.Header.Set("API_KEY", "token_basic")
		
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro na requisição %d com token básico: %v", i+1, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Requisição %d com token básico deveria ser permitida, status: %d", i+1, resp.StatusCode)
		}
	}
	
	// Terceira requisição com token básico deve ser bloqueada
	req, _ := http.NewRequest("GET", baseURL+"/teste", nil)
	req.Header.Set("API_KEY", "token_basic")
	
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Erro na terceira requisição com token básico: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Terceira requisição com token básico deveria ser bloqueada, status: %d", resp.StatusCode)
	}
	
	// Teste com token VIP (limite 10) - deve permitir mais requisições
	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", baseURL+"/teste", nil)
		req.Header.Set("API_KEY", "token_vip")
		
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro na requisição %d com token VIP: %v", i+1, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Requisição %d com token VIP deveria ser permitida, status: %d", i+1, resp.StatusCode)
		}
	}
}

func testeConcorrenciaIP(t *testing.T, baseURL string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	sucessos := 0
	bloqueios := 0
	
	// Fazer 10 requisições concorrentes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			resp, err := client.Get(baseURL + "/teste")
			if err != nil {
				t.Errorf("Erro na requisição concorrente: %v", err)
				return
			}
			defer resp.Body.Close()
			
			mu.Lock()
			if resp.StatusCode == http.StatusOK {
				sucessos++
			} else if resp.StatusCode == http.StatusTooManyRequests {
				bloqueios++
			}
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	// Verificar que o número total de respostas está correto
	total := sucessos + bloqueios
	if total != 10 {
		t.Errorf("Esperado 10 respostas totais, mas obteve %d", total)
	}
	
	// Deve ter no máximo 3 sucessos (respeitando o limite)
	if sucessos > 3 {
		t.Errorf("Sucessos não devem exceder o limite de 3, mas obteve %d", sucessos)
	}
	
	// Deve ter pelo menos alguns bloqueios
	if bloqueios == 0 {
		t.Error("Deveria haver pelo menos alguns bloqueios em requisições concorrentes")
	}
	
	t.Logf("Resultado do teste de concorrência: %d sucessos, %d bloqueios", sucessos, bloqueios)
}

func testeDiferentesIPs(t *testing.T, baseURL string) {
	// Simular diferentes IPs usando headers X-Forwarded-For
	client := &http.Client{Timeout: 5 * time.Second}
	
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	
	for _, ip := range ips {
		// Cada IP deve conseguir fazer 3 requisições
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest("GET", baseURL+"/teste", nil)
			req.Header.Set("X-Forwarded-For", ip)
			
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Erro na requisição %d para IP %s: %v", i+1, ip, err)
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição %d para IP %s deveria ser permitida, status: %d", i+1, ip, resp.StatusCode)
			}
		}
		
		// Quarta requisição deve ser bloqueada
		req, _ := http.NewRequest("GET", baseURL+"/teste", nil)
		req.Header.Set("X-Forwarded-For", ip)
		
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Erro na quarta requisição para IP %s: %v", ip, err)
		}
		defer resp.Body.Close()
		
		if resp.StatusCode != http.StatusTooManyRequests {
			t.Errorf("Quarta requisição para IP %s deveria ser bloqueada, status: %d", ip, resp.StatusCode)
		}
	}
}

func TestIntegracao_ConfiguracaoCompleta(t *testing.T) {
	// Teste de carregamento de configuração
	cfg, err := config.CarregarConfig()
	if err != nil {
		t.Fatalf("Erro ao carregar configuração: %v", err)
	}
	
	// Verificar valores padrão
	if cfg.PortaServidor != 8080 {
		t.Errorf("Porta padrão incorreta: esperado 8080, obtido %d", cfg.PortaServidor)
	}
	
	if cfg.LimiteIPPorSegundo != 10 {
		t.Errorf("Limite IP padrão incorreto: esperado 10, obtido %d", cfg.LimiteIPPorSegundo)
	}
	
	if cfg.LimiteTokenPorSegundo != 100 {
		t.Errorf("Limite token padrão incorreto: esperado 100, obtido %d", cfg.LimiteTokenPorSegundo)
	}
	
	if cfg.TempoBloqueioIP != 5*time.Minute {
		t.Errorf("Tempo bloqueio IP incorreto: esperado 5m, obtido %v", cfg.TempoBloqueioIP)
	}
}

// Benchmark de performance
func BenchmarkRateLimiter_RequisicoesConcorrentes(b *testing.B) {
	servidor := criarServidorTeste()
	defer servidor.Close()
	
	client := &http.Client{Timeout: 5 * time.Second}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Simular diferentes IPs para evitar rate limiting
			req, _ := http.NewRequest("GET", servidor.URL+"/benchmark", nil)
			req.Header.Set("X-Forwarded-For", fmt.Sprintf("192.168.1.%d", (i%250)+1))
			
			resp, err := client.Do(req)
			if err != nil {
				b.Errorf("Erro na requisição: %v", err)
				continue
			}
			resp.Body.Close()
			i++
		}
	})
}