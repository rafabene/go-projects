## Como executar o projeto

### 1. Clone o repositório

```bash
git clone https://github.com/rafabene/go-projects.git
cd go-projects/cepweather
```

Antes de executar os demais passos, crie o arquivo `.env` copiando o template e preenchendo sua chave da WeatherAPI:

```bash
cp configs/.env.template configs/.env
# Edite configs/.env e coloque sua chave da WeatherAPI
# que pode ser obtida em https://www.weatherapi.com/
```

### 2. Execute os testes

```bash
go test ./tests
```

### 3. Build e execute com Docker Compose

```bash
docker-compose -f docker/docker-compose.yml up --build
```

O serviço estará disponível em `http://localhost:8080`.

### 4. Teste a API com curl

```bash
curl -X POST http://localhost:8080/api/v1/weather \
  -H "Content-Type: application/json" \
  -d '{"cep":"75110430"}'
```

Você deve receber uma resposta com os dados do clima para o CEP informado.
