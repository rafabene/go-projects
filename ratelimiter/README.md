# Rate Limiter em Go

Um sistema de controle de taxa (rate limiter) ultra-simplificado em Go que permite limitar o número de requisições por segundo com base em endereço IP ou token de acesso.

## 📋 Características

- **Limitação por IP**: Controla requisições baseado no endereço IP do cliente
- **Limitação por Token**: Controla requisições baseado em tokens de acesso (header `API_KEY`)
- **Prioridade de Token**: Configurações de token sobrepõem as de IP
- **Tempos de Bloqueio**: Configura tempo de bloqueio quando limite é excedido
- **Ultra-Simples**: Implementação direta sem abstrações desnecessárias
- **Thread-Safe**: Usa `sync.Map` nativo do Go para concorrência
- **Middleware HTTP**: Integração fácil com servidores web Go
- **Configuração Flexível**: Via variáveis de ambiente ou arquivo `.env`

## 🏗️ Arquitetura Simplificada

```
├── cmd/servidor/           # Aplicação principal
├── internal/
│   ├── config/            # Configurações
│   └── middleware/        # Rate limiter + middleware HTTP
├── test/                  # Testes de integração
├── .env                   # Configurações
├── Dockerfile            # Container Docker
└── docker-compose.yml   # Orquestração
```

## 📦 Dependências

- **Go 1.21+**: Linguagem de programação
- **[github.com/joho/godotenv](https://github.com/joho/godotenv)**: Carregamento de arquivos `.env`

## 🚀 Como Usar

### Executar Localmente

1. **Clone e configure**:
   ```bash
   git clone https://github.com/rafabene/go-projects/ratelimiter.git
   cd ratelimiter
   ```

2. **Instale as dependências**:
   ```bash
   go mod download
   ```

3. **Configure as variáveis (opcional)**:
   ```bash
   # Edite o arquivo .env ou exporte variáveis de ambiente
   export LIMITE_IP_POR_SEGUNDO=10
   export TEMPO_BLOQUEIO_IP=300
   export LIMITE_TOKEN_POR_SEGUNDO=100
   export TEMPO_BLOQUEIO_TOKEN=300
   export PORTA_SERVIDOR=8080
   ```

4. **Execute o servidor**:
   ```bash
   go run cmd/servidor/main.go
   ```

### Executar com Docker

```bash
# Construir e executar
docker-compose up --build

# Ou apenas executar se já foi construído
docker-compose up
```

## ⚙️ Configuração

O projeto utiliza o [`github.com/joho/godotenv`](https://github.com/joho/godotenv) para carregamento de configurações.

### Prioridade de Configuração

1. **Variáveis de ambiente do sistema** (maior prioridade)
2. **Arquivo `.env`** (configurações do projeto)
3. **Valores padrão do código** (menor prioridade)

### Variáveis de Ambiente

| Variável | Descrição | Valor Padrão |
|----------|-----------|--------------|
| `PORTA_SERVIDOR` | Porta do servidor HTTP | `8080` |
| `LIMITE_IP_POR_SEGUNDO` | Requisições permitidas por IP/segundo | `10` |
| `TEMPO_BLOQUEIO_IP` | Tempo de bloqueio do IP (segundos) | `300` |
| `LIMITE_TOKEN_POR_SEGUNDO` | Requisições permitidas por token/segundo | `100` |
| `TEMPO_BLOQUEIO_TOKEN` | Tempo de bloqueio do token (segundos) | `300` |

### Tokens Personalizados

Para configurar limites específicos por token, use variáveis no formato:
```bash
TOKEN_LIMITE_<nome_do_token>=<limite>
```

Exemplo no `.env`:
```bash
TOKEN_LIMITE_abc123=200
TOKEN_LIMITE_vip_user=500
TOKEN_LIMITE_basic_user=50
```

### Arquivo de Configuração

**`.env`** (configuração do projeto):
```bash
# Configuração do Rate Limiter
PORTA_SERVIDOR=8080
LIMITE_IP_POR_SEGUNDO=10
TEMPO_BLOQUEIO_IP=300
LIMITE_TOKEN_POR_SEGUNDO=100
TEMPO_BLOQUEIO_TOKEN=300

# Tokens personalizados (opcional)
TOKEN_LIMITE_abc123=200
TOKEN_LIMITE_vip_user=500
TOKEN_LIMITE_basic_user=50
```

## 🔌 Endpoints

- `GET /` - Endpoint principal de teste
- `GET /status` - Status do servidor
- `GET /teste` - Endpoint para testes de carga

## 🧪 Testando o Rate Limiter

### Teste Básico por IP

```bash
# Fazer múltiplas requisições rápidas
for i in {1..15}; do
  curl -w "Status: %{http_code}\n" http://localhost:8080/
done
```

### Teste com Token

```bash
# Com token personalizado
curl -H "API_KEY: abc123" http://localhost:8080/

# Sem token (usa limite do IP)
curl http://localhost:8080/
```

### Teste de Carga com Apache Bench

```bash
# 100 requisições, 10 concorrentes
ab -n 100 -c 10 http://localhost:8080/teste

# Com token
ab -n 100 -c 10 -H "API_KEY: vip_user" http://localhost:8080/teste
```

## 📊 Comportamento do Rate Limiter

### Limitação por IP
- Cada IP único tem seu próprio contador
- Quando excede o limite: retorna status `429`
- Após o tempo de bloqueio: contador é resetado

### Limitação por Token
- Token no header `API_KEY` tem prioridade sobre IP
- Tokens personalizados podem ter limites diferentes
- Mesmo IP com token diferente = contadores separados

### Resposta de Erro (429)
```json
{
  "erro": "you have reached the maximum number of requests or actions allowed within a certain time frame",
  "codigo": 429,
  "detalhes": "Limite excedido para IP. Tente novamente em 5m0s"
}
```

## 🔧 Executar Testes

```bash
# Testes do middleware (principal)
go test ./internal/middleware -v

# Testes de configuração
go test ./internal/config -v

# Testes de integração
go test ./test -v

# Todos os testes
go test ./... -v

# Benchmark
go test ./... -bench=. -benchmem
```

## 🐳 Docker

### Construir Imagem
```bash
docker build -t ratelimiter .
```

### Executar Container
```bash
docker run -p 8080:8080 \
  -e LIMITE_IP_POR_SEGUNDO=5 \
  -e TOKEN_LIMITE_vip=100 \
  ratelimiter
```

### Docker Compose
```bash
# Subir todos os serviços
docker-compose up -d

# Ver logs
docker-compose logs -f

# Parar serviços
docker-compose down
```

## 📈 Monitoramento

### Health Check
```bash
curl http://localhost:8080/status
```

### Logs
O servidor exibe logs detalhados incluindo:
- Configurações carregadas
- Requisições bloqueadas/permitidas
- Erros de sistema

## 🤝 Exemplo de Uso Personalizado

### Middleware Simples
```go
package main

import (
    "net/http"
    "time"
    
    "github.com/rafabene/go-projects/ratelimiter/internal/middleware"
)

func main() {
    // Configuração simples
    config := &middleware.ConfigRateLimiter{
        LimiteIPPorSegundo:    5,
        TempoBloqueioIP:       time.Minute,
        LimiteTokenPorSegundo: 50,
        TempoBloqueioToken:    5 * time.Minute,
        TokensPersonalizados:  map[string]int{
            "vip": 100,
        },
    }
    
    // Criar rate limiter (sem dependências externas!)
    rateLimiter := middleware.NovoRateLimiter(config)
    
    // Aplicar a todas as rotas
    http.Handle("/api/", rateLimiter.Middleware(meuHandler))
    http.ListenAndServe(":8080", nil)
}
```

## ⚡ Performance

- **Concorrência**: Thread-safe usando `sync.Map` nativo
- **Simplicidade**: Zero abstrações desnecessárias
- **Latência**: < 1ms por verificação de rate limit
- **Throughput**: > 10,000 req/s em hardware moderno
- **Memória**: Baixo uso, sem limpeza automática complexa

## 🎯 Arquitetura Ultra-Simples

### Componentes Principais

1. **`middleware.RateLimiter`**: Struct principal com `sync.Map` interno
2. **`middleware.NovoRateLimiter(config)`**: Construtor simples
3. **`config.CarregarConfig()`**: Carregamento de configurações
4. **Endpoints HTTP**: Três endpoints de teste

### Fluxo de Funcionamento

```
Requisição HTTP → Middleware → sync.Map → Permit/Deny → Response
```

## 🐛 Solução de Problemas

### Servidor não inicia
- Verifique se a porta 8080 está livre
- Confirme as permissões do arquivo `.env`

### Rate limiting não funciona
- Verifique se o middleware está sendo aplicado
- Confirme as configurações de limite

### Testes falhando
- Aguarde reset dos contadores entre testes
- Execute testes isoladamente se necessário

## 🏆 Vantagens da Simplificação

✅ **Zero abstrações desnecessárias**  
✅ **Uma única struct principal**  
✅ **sync.Map nativo para thread safety**  
✅ **Construtor sem dependências**  
✅ **Código fácil de entender e manter**  
✅ **Performance máxima**  

## 📝 Licença

Este projeto é desenvolvido para fins educacionais e demonstrativos.

---

**Desenvolvido em Go com ❤️ - Versão Ultra-Simplificada**