# Stress Test CLI

Uma ferramenta CLI desenvolvida em Go para realizar testes de carga em serviços web. Permite especificar URL, número de requests e nível de concorrência, fornecendo relatórios detalhados sobre o desempenho.

## 🚀 Funcionalidades

- ⚡ Testes de carga com controle de concorrência
- 📊 Relatórios detalhados com métricas de desempenho
- 🐳 Execução via Docker
- 🎯 Interface CLI amigável com Cobra
- 📈 Distribuição de códigos de status HTTP
- ⏱️ Medição de tempo de resposta e throughput

## 📋 Pré-requisitos

### Para execução local:
- Go 1.23 ou superior
- Acesso à internet para testes

### Para execução via Docker:
- Docker instalado

## 🛠️ Instalação e Execução

### Método 1: Execução Local

1. **Clone o repositório:**
```bash
git clone https://github.com/rafabene/go-projects.git
cd go-projects/stress-test
```

2. **Instale as dependências:**
```bash
go mod tidy
```

3. **Compile a aplicação:**
```bash
go build -o stress-test ./cmd/main
```

4. **Execute os testes:**
```bash
./stress-test --url=https://httpbin.org/status/200 --requests=100 --concurrency=10
```

### Método 2: Execução via Docker

1. **Construa a imagem Docker:**
```bash
docker build -t stress-test .
```

2. **Execute o container:**
```bash
docker run stress-test --url=https://httpbin.org/status/200 --requests=100 --concurrency=10
```

## 📖 Uso

### Sintaxe

```bash
stress-test [flags]
```

### Flags Disponíveis

| Flag | Flag Curta | Descrição | Obrigatório | Padrão |
|------|------------|-----------|-------------|---------|
| `--url` | `-u` | URL do serviço a ser testado | ✅ | - |
| `--requests` | `-r` | Número total de requests | ✅ | - |
| `--concurrency` | `-c` | Número de chamadas simultâneas | ❌ | 1 |
| `--help` | `-h` | Exibe informações de ajuda | ❌ | - |

### Exemplos de Uso

**Teste básico:**
```bash
./stress-test -u https://google.com -r 50 -c 5
```

**Teste com alta concorrência:**
```bash
./stress-test --url=https://api.exemplo.com/health --requests=1000 --concurrency=50
```

**Teste via Docker:**
```bash
docker run stress-test -u https://httpbin.org/status/200 -r 200 -c 20
```

**Visualizar ajuda:**
```bash
./stress-test --help
```

## 📊 Interpretando o Relatório

O relatório gerado inclui as seguintes métricas:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 RELATÓRIO DO TESTE DE CARGA
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
⏱️  Tempo total gasto: 5.123s
📨 Total de requests realizados: 100
✅ Requests com status 200: 95
❌ Requests com erro: 2

📈 Distribuição de códigos de status:
   200: 95 requests (95.0%)
   404: 3 requests (3.0%)

🚀 Requests por segundo: 19.52 req/s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### Métricas Explicadas

- **Tempo total gasto**: Duração total do teste
- **Total de requests realizados**: Número de requisições executadas
- **Requests com status 200**: Requisições bem-sucedidas
- **Requests com erro**: Requisições que falharam (timeout, erro de rede, etc.)
- **Distribuição de códigos de status**: Breakdown detalhado dos códigos HTTP retornados
- **Requests por segundo**: Taxa de throughput (RPS)

## 🏗️ Estrutura do Projeto

```
├── cmd/main/
│   └── main.go              # Ponto de entrada da aplicação
├── pkg/stresstest/
│   ├── config.go            # Configuração e validação
│   ├── types.go             # Definições de tipos
│   ├── runner.go            # Lógica de execução dos testes
│   └── reporter.go          # Geração de relatórios
├── Dockerfile               # Configuração Docker
├── go.mod                   # Dependências Go
├── go.sum                   # Checksums das dependências
└── README.md               # Este arquivo
```

## 🔧 Desenvolvimento

### Executar em modo desenvolvimento:
```bash
go run ./cmd/main --url=https://httpbin.org/status/200 --requests=10 --concurrency=3
```

### Executar testes:
```bash
go test ./...
```

### Compilar para diferentes plataformas:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o stress-test-linux ./cmd/main

# Windows
GOOS=windows GOARCH=amd64 go build -o stress-test.exe ./cmd/main

# macOS
GOOS=darwin GOARCH=amd64 go build -o stress-test-mac ./cmd/main
```

## ⚠️ Considerações Importantes

1. **Responsabilidade**: Use esta ferramenta apenas em serviços que você possui ou tem permissão para testar
2. **Rate Limiting**: Alguns serviços podem ter limitação de taxa. Ajuste a concorrência adequadamente
3. **Recursos do Sistema**: Testes com alta concorrência podem consumir muitos recursos de rede e CPU
4. **Timeout**: Cada request tem timeout de 30 segundos

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença Apache 2.0. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🐛 Reportando Bugs

Se encontrar algum bug ou tiver sugestões, por favor abra uma [issue](https://github.com/rafabene/go-projects/issues).

---

**Desenvolvido com ❤️ em Go**