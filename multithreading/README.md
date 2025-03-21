# Projeto de Multithreading em Go

Este projeto demonstra a implementação de conceitos de multithreading em Go, utilizando goroutines e canais para processamento concorrente.

## Estrutura do Projeto

```
multithreading/
├── cmd/            # Diretório contendo os executáveis principais
├── config/         # Configurações e utilitários
├── config.yaml     # Arquivo de configuração
├── go.mod          # Arquivo de dependências Go
└── go.sum          # Checksums das dependências
```

## Requisitos

- Go 1.23.6 ou superior
- Dependências listadas no `go.mod`

## Instalação

1. Clone o repositório:

```bash
git clone https://github.com/rafabene/go-projects.git
cd go-projects/multithreading
```

2. Instale as dependências:

```bash
go mod download
```

## Como Usar

Para executar o projeto:

```bash
go run cmd/fast.go
```

## Configuração

O projeto utiliza um arquivo `config.yaml` para configurações. Você pode modificar as configurações editando este arquivo.

## Funcionalidades

- Implementação de goroutines para processamento concorrente
- Uso de canais para comunicação entre goroutines
- Configuração via arquivo YAML
- Gerenciamento de configurações com Viper

## Contribuição

Contribuições são bem-vindas! Por favor, sinta-se à vontade para submeter pull requests.

## Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
