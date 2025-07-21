## Como executar o projeto

## CEPWeather API

API REST em Go para consulta de clima atual por CEP brasileiro, integrando com a WeatherAPI. Ideal para uso em sistemas que precisam de informações meteorológicas baseadas em localização.

---

### Funcionalidades

- Consulta o clima atual para um CEP informado.
- Validação de formato e tamanho do CEP.
- Retorno estruturado em JSON.
- Log detalhado das requisições e respostas.
- Pronto para deploy local, Docker e Google Cloud Run.

---

### Estrutura do Projeto

```
serviceB/
├── cmd/                # Ponto de entrada da aplicação (main.go)
├── internal/
│   ├── handlers/       # Handlers HTTP
│   ├── models/         # Structs e tipos
│   └── services/       # Integração com APIs externas
├── configs/            # Arquivos de configuração (.env)
├── tests/              # Testes automatizados
├── docker/             # Arquivos para execução em containers
└── README.md           # Documentação
```

---

### Parâmetros de Entrada

**POST /api/v1/weather**

Body (JSON):

```json
{
  "cep": "75110430"
}
```

---

### Parâmetros de Saída

Resposta 200 OK:

```json
{
  "temp_c": 25.0,
  "temp_f": 77.0,
  "temp_k": 298.15
}
```

Erros possíveis:

- 400 Bad Request: JSON inválido
- 422 Unprocessable Entity: CEP inválido
- 500 Internal Server Error: Falha na integração externa

---

### Como executar o projeto

#### 1. Clone o repositório

```bash
git clone https://github.com/rafabene/go-projects.git
cd go-projects/cepweather-otel/serviceB
```

#### 2. Configure a chave da WeatherAPI

```bash
cp configs/.env.template configs/.env
# Edite configs/.env e coloque sua chave da WeatherAPI
# que pode ser obtida em https://www.weatherapi.com/
```

#### 3. Execute os testes

```bash
go test ./tests
```

#### 4. Build e execute

```bash
go build -o serviceb cmd/main.go
./serviceb
```

O serviço estará disponível em `http://localhost:8081`.

#### 5. Teste a API com curl

```bash
curl -X POST http://localhost:8081/api/v1/weather \
  -H "Content-Type: application/json" \
  -d '{"cep":"75110430"}'
```

Você deve receber uma resposta com os dados do clima para o CEP informado.

---

### Informações técnicas

- Linguagem: Go 1.21+
- Framework: net/http
- Integração: WeatherAPI
- Testes: Go testing/httptest
- Docker-ready
- Variáveis de ambiente via `.env`

---

### Autor

Rafael Benevides
https://github.com/rafabene
