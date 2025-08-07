// Package middleware implementa middleware HTTP para rate limiting.
//
// Este pacote fornece um middleware que pode ser facilmente integrado
// a qualquer aplicação web Go, oferecendo controle de taxa baseado em
// IP ou token de acesso com configuração flexível.
package middleware

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ConfigRateLimiter contém todas as configurações necessárias para o rate limiter.
//
// A configuração permite definir limites diferentes para IPs e tokens,
// com a possibilidade de configurar tokens específicos com limites personalizados.
type ConfigRateLimiter struct {
	LimiteIPPorSegundo    int                // Limite de requisições por IP por segundo
	TempoBloqueioIP       time.Duration      // Tempo de bloqueio quando IP excede limite
	LimiteTokenPorSegundo int                // Limite padrão de requisições por token por segundo
	TempoBloqueioToken    time.Duration      // Tempo de bloqueio quando token excede limite
	TokensPersonalizados  map[string]int     // Limites específicos por token (chave: token, valor: limite)
}

// InformacaoLimite contém informações simples sobre um limitador.
type InformacaoLimite struct {
	Contador  int       // Número de requisições no período atual
	UltimaVez time.Time // Timestamp da última requisição
}

// RateLimiter é o middleware HTTP que implementa controle de taxa.
//
// O middleware funciona interceptando todas as requisições HTTP e
// aplicando as regras de rate limiting antes de passar para o próximo handler.
//
// Prioridade de verificação:
//  1. Se há token API_KEY -> aplica limite do token (sobrepõe IP)
//  2. Se não há token -> aplica limite por IP
type RateLimiter struct {
	limites sync.Map          // Thread-safe map para armazenar limitadores
	config  *ConfigRateLimiter // Configurações de limite e bloqueio
}

// NovoRateLimiter cria uma nova instância do middleware rate limiter.
//
// Parâmetros:
//   - config: configurações de limites e tempos de bloqueio
//
// Retorna um middleware pronto para ser usado com qualquer router HTTP.
func NovoRateLimiter(config *ConfigRateLimiter) *RateLimiter {
	return &RateLimiter{
		config: config,
	}
}

// RespostaErro representa a estrutura JSON retornada quando o limite é excedido.
//
// Segue o padrão HTTP 429 (Too Many Requests) com informações detalhadas
// sobre o erro para facilitar o debugging e orientar o cliente.
type RespostaErro struct {
	Erro     string `json:"erro"`               // Mensagem de erro padrão
	Codigo   int    `json:"codigo"`             // Código HTTP (sempre 429)
	Detalhes string `json:"detalhes,omitempty"` // Informações adicionais sobre o bloqueio
}

// Middleware retorna o middleware HTTP que implementa rate limiting.
//
// O middleware intercepta todas as requisições HTTP e aplica as regras de
// controle de taxa antes de passar a requisição para o handler seguinte.
//
// Fluxo de processamento:
//  1. Extrai o IP real do cliente (considerando proxies)
//  2. Verifica se há token API_KEY no header
//  3. Se há token -> aplica limite por token (prioridade)
//  4. Se não há token -> aplica limite por IP
//  5. Se bloqueado -> retorna HTTP 429 com detalhes
//  6. Se permitido -> continua para o próximo handler
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrai o IP real do cliente, considerando headers de proxy
		ip := rl.extrairIP(r)
		
		// Extrai token de acesso do header API_KEY
		token := r.Header.Get("API_KEY")
		
		// Token tem prioridade sobre IP - se existe token, usa limite de token
		if token != "" {
			permitido, tempoEspera := rl.verificarLimiteToken(token)
			if !permitido {
				rl.enviarErroLimite(w, tempoEspera, "token")
				return
			}
		} else {
			// Sem token - aplica limitação por IP
			permitido, tempoEspera := rl.verificarLimiteIP(ip)
			if !permitido {
				rl.enviarErroLimite(w, tempoEspera, "IP")
				return
			}
		}
		
		// Requisição permitida - continua para o próximo handler
		next.ServeHTTP(w, r)
	})
}

// extrairIP extrai o endereço IP real do cliente considerando proxies e load balancers.
//
// A função segue a ordem de prioridade recomendada para detectar o IP real:
//  1. X-Forwarded-For (primeiro IP da lista em caso de múltiplos proxies)
//  2. X-Real-IP (IP definido por proxies reversos como nginx)  
//  3. RemoteAddr (conexão direta quando não há proxies)
//
// Esta implementação é importante para garantir que o rate limiting
// funcione corretamente mesmo atrás de proxies, CDNs ou load balancers.
func (rl *RateLimiter) extrairIP(r *http.Request) string {
	// X-Forwarded-For: lista de IPs separados por vírgula (cliente, proxy1, proxy2, ...)
	// O primeiro IP da lista é o cliente original
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	
	// X-Real-IP: IP do cliente definido por proxy reverso (nginx, etc.)
	if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
		return xRealIP
	}
	
	// RemoteAddr: IP da conexão direta (sem proxy)
	// Formato: "IP:porta", precisamos extrair apenas o IP
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// Se não conseguir dividir, retorna o valor completo
		return r.RemoteAddr
	}
	
	return ip
}

// verificarLimiteToken verifica se um token de acesso excedeu seu limite de requisições.
//
// A função aplica limites específicos por token, onde tokens personalizados
// podem ter limites diferentes do padrão configurado.
//
// Processo:
//  1. Verifica se o token tem limite personalizado configurado
//  2. Usa o limite específico ou o padrão se não configurado
//  3. Cria chave única para o token ("token:abc123")
//  4. Consulta a estratégia para verificar se pode fazer a requisição
//  5. Se excedeu, aplica tempo de bloqueio configurado
//
// Em caso de erro na estratégia, permite a requisição (fail-open) para
// evitar quebrar o serviço por problemas de infraestrutura.
func (rl *RateLimiter) verificarLimiteToken(token string) (bool, time.Duration) {
	// Determina o limite aplicável para este token
	limite := rl.config.LimiteTokenPorSegundo
	if limitePersonalizado, existe := rl.config.TokensPersonalizados[token]; existe {
		limite = limitePersonalizado
	}
	
	// Cria chave única para identificar este token
	chave := fmt.Sprintf("token:%s", token)
	return rl.permitirRequisicao(chave, limite, time.Second)
}

// verificarLimiteIP verifica se um endereço IP excedeu seu limite de requisições.
//
// Aplica o limite global configurado para todos os IPs, sem diferenciação.
// É usado quando não há token de acesso presente na requisição.
//
// O processo é similar ao de tokens, mas mais simples pois não há
// limites personalizados por IP (todos usam o mesmo limite global).
func (rl *RateLimiter) verificarLimiteIP(ip string) (bool, time.Duration) {
	// Cria chave única para identificar este IP
	chave := fmt.Sprintf("ip:%s", ip)
	return rl.permitirRequisicao(chave, rl.config.LimiteIPPorSegundo, time.Second)
}

// permitirRequisicao implementa rate limiting simples.
func (rl *RateLimiter) permitirRequisicao(chave string, limite int, janelaTempo time.Duration) (bool, time.Duration) {
	agora := time.Now()
	
	// Carrega ou cria nova informação de limite
	value, _ := rl.limites.LoadOrStore(chave, &InformacaoLimite{
		Contador:  0,
		UltimaVez: agora,
	})
	
	info := value.(*InformacaoLimite)
	
	// Se passou da janela de tempo, reseta
	if agora.Sub(info.UltimaVez) >= janelaTempo {
		info.Contador = 1
		info.UltimaVez = agora
		return true, 0
	}
	
	// Incrementa contador
	info.Contador++
	info.UltimaVez = agora
	
	// Verifica se excedeu limite
	if info.Contador > limite {
		return false, janelaTempo
	}
	
	return true, 0
}

// enviarErroLimite envia resposta HTTP 429 quando o limite de taxa é excedido.
//
// A função configura os headers HTTP apropriados conforme RFC 6585 e
// retorna uma resposta JSON estruturada com informações sobre o erro.
//
// Headers configurados:
//   - Content-Type: application/json
//   - Retry-After: tempo em segundos até poder tentar novamente
//   - Status: 429 Too Many Requests
//
// A resposta inclui detalhes em português para facilitar o debugging.
func (rl *RateLimiter) enviarErroLimite(w http.ResponseWriter, tempoEspera time.Duration, tipo string) {
	// Configura headers HTTP apropriados para rate limiting
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", strconv.Itoa(int(tempoEspera.Seconds())))
	w.WriteHeader(http.StatusTooManyRequests)
	
	// Cria resposta estruturada conforme especificação
	resposta := RespostaErro{
		Erro:     "you have reached the maximum number of requests or actions allowed within a certain time frame",
		Codigo:   http.StatusTooManyRequests,
		Detalhes: fmt.Sprintf("Limite excedido para %s. Tente novamente em %v", tipo, tempoEspera),
	}
	
	// Envia resposta JSON (ignora erro de encoding pois é estrutura simples)
	json.NewEncoder(w).Encode(resposta)
}