package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_LimitacaoIP(t *testing.T) {
	// Configurar rate limiter
	
	config := &ConfigRateLimiter{
		LimiteIPPorSegundo: 3,
		TempoBloqueioIP:    5 * time.Second,
		LimiteTokenPorSegundo: 100,
		TempoBloqueioToken: 5 * time.Second,
		TokensPersonalizados: make(map[string]int),
	}
	
	rateLimiter := NovoRateLimiter(config)
	
	// Handler de teste
	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sucesso"))
	}))
	
	// Testar requisições permitidas
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Requisição %d deveria ser permitida, mas retornou %d", i+1, rr.Code)
		}
	}
	
	// Quarta requisição deve ser bloqueada
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Quarta requisição deveria ser bloqueada, mas retornou %d", rr.Code)
	}
	
	// Verificar se header Retry-After está presente
	retryAfter := rr.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Error("Header Retry-After deveria estar presente")
	}
}

func TestRateLimiter_LimitacaoToken(t *testing.T) {
	// Configurar rate limiter
	
	config := &ConfigRateLimiter{
		LimiteIPPorSegundo: 100,
		TempoBloqueioIP:    5 * time.Second,
		LimiteTokenPorSegundo: 2,
		TempoBloqueioToken: 5 * time.Second,
		TokensPersonalizados: make(map[string]int),
	}
	
	rateLimiter := NovoRateLimiter(config)
	
	// Handler de teste
	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sucesso"))
	}))
	
	// Testar requisições com token permitidas
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "token123")
		req.RemoteAddr = "192.168.1.1:12345"
		
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Requisição %d com token deveria ser permitida, mas retornou %d", i+1, rr.Code)
		}
	}
	
	// Terceira requisição deve ser bloqueada
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("API_KEY", "token123")
	req.RemoteAddr = "192.168.1.1:12345"
	
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Terceira requisição com token deveria ser bloqueada, mas retornou %d", rr.Code)
	}
}

func TestRateLimiter_TokenPersonalizado(t *testing.T) {
	// Configurar rate limiter com token personalizado
	
	tokensPersonalizados := make(map[string]int)
	tokensPersonalizados["token_vip"] = 10
	
	config := &ConfigRateLimiter{
		LimiteIPPorSegundo: 3,
		TempoBloqueioIP:    5 * time.Second,
		LimiteTokenPorSegundo: 5,
		TempoBloqueioToken: 5 * time.Second,
		TokensPersonalizados: tokensPersonalizados,
	}
	
	rateLimiter := NovoRateLimiter(config)
	
	// Handler de teste
	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sucesso"))
	}))
	
	// Testar que token VIP pode fazer mais requisições
	for i := 0; i < 8; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("API_KEY", "token_vip")
		req.RemoteAddr = "192.168.1.1:12345"
		
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		
		if rr.Code != http.StatusOK {
			t.Errorf("Requisição %d com token VIP deveria ser permitida, mas retornou %d", i+1, rr.Code)
		}
	}
}

func TestRateLimiter_PrioridadeTokenSobreIP(t *testing.T) {
	// Configurar rate limiter
	
	config := &ConfigRateLimiter{
		LimiteIPPorSegundo: 1, // Limite muito baixo por IP
		TempoBloqueioIP:    5 * time.Second,
		LimiteTokenPorSegundo: 5, // Limite mais alto por token
		TempoBloqueioToken: 5 * time.Second,
		TokensPersonalizados: make(map[string]int),
	}
	
	rateLimiter := NovoRateLimiter(config)
	
	// Handler de teste
	handler := rateLimiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("sucesso"))
	}))
	
	// Primeira requisição sem token - deve ser permitida
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "192.168.1.1:12345"
	rr1 := httptest.NewRecorder()
	handler.ServeHTTP(rr1, req1)
	
	if rr1.Code != http.StatusOK {
		t.Errorf("Primeira requisição sem token deveria ser permitida, mas retornou %d", rr1.Code)
	}
	
	// Segunda requisição sem token - deve ser bloqueada (limite IP = 1)
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "192.168.1.1:12345"
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	
	if rr2.Code != http.StatusTooManyRequests {
		t.Errorf("Segunda requisição sem token deveria ser bloqueada, mas retornou %d", rr2.Code)
	}
	
	// Requisição com token do mesmo IP - deve ser permitida (prioridade do token)
	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("API_KEY", "token123")
	req3.RemoteAddr = "192.168.1.1:12345"
	rr3 := httptest.NewRecorder()
	handler.ServeHTTP(rr3, req3)
	
	if rr3.Code != http.StatusOK {
		t.Errorf("Requisição com token deveria ser permitida (prioridade), mas retornou %d", rr3.Code)
	}
}

func TestRateLimiter_ExtrairIP(t *testing.T) {
	
	config := &ConfigRateLimiter{
		LimiteIPPorSegundo: 10,
		TempoBloqueioIP:    5 * time.Second,
		LimiteTokenPorSegundo: 10,
		TempoBloqueioToken: 5 * time.Second,
		TokensPersonalizados: make(map[string]int),
	}
	
	rateLimiter := NovoRateLimiter(config)
	
	testes := []struct {
		nome     string
		headers  map[string]string
		remoteAddr string
		esperado string
	}{
		{
			nome:       "X-Forwarded-For",
			headers:    map[string]string{"X-Forwarded-For": "203.0.113.1, 70.41.3.18"},
			remoteAddr: "192.168.1.1:12345",
			esperado:   "203.0.113.1",
		},
		{
			nome:       "X-Real-IP",
			headers:    map[string]string{"X-Real-IP": "203.0.113.2"},
			remoteAddr: "192.168.1.1:12345",
			esperado:   "203.0.113.2",
		},
		{
			nome:       "RemoteAddr",
			headers:    map[string]string{},
			remoteAddr: "192.168.1.1:12345",
			esperado:   "192.168.1.1",
		},
	}
	
	for _, teste := range testes {
		t.Run(teste.nome, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = teste.remoteAddr
			
			for chave, valor := range teste.headers {
				req.Header.Set(chave, valor)
			}
			
			ip := rateLimiter.extrairIP(req)
			if ip != teste.esperado {
				t.Errorf("IP extraído incorreto. Esperado: %s, Obtido: %s", teste.esperado, ip)
			}
		})
	}
}