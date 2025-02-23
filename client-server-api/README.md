# 📌 Projeto: Cotação do Dólar

Este projeto implementa um servidor e um cliente para obter a cotação do dólar e armazenar os dados utilizando **Go** e **contextos para controle de timeout**.

## 🚀 Como Executar o Projeto

### 🔹 **1. Iniciar o Servidor**

```sh
cd server
go run server.go
```

### 🔹 **2. Iniciar o Cliente**

```sh
cd client
go run client.go
```

## 🛠️ Tecnologias Utilizadas

- **Go (Golang)**
- **net/http**
- **context (para timeouts)**
- **cURL para testes**

## 🔄 Fluxo do Sistema

1️⃣ **Cliente (`client.go`)** solicita a cotação do dólar ao servidor.

2️⃣ **Servidor (`server.go`)** busca a cotação da API externa e futuramente armazenará no banco SQLite.

3️⃣ **Servidor responde ao cliente** com o valor do **bid**.

4️⃣ **Cliente salva o valor da cotação** no arquivo `cotacoes.txt`.

---

🚀 **Codado por Rafael Benevides** em 23/02/2025 🚀
