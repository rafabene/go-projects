package main

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/rafabene/client-server-api/common"
)

func main() {
	cotacao := obterCotacao()
	println("Cotação do Dólar: ", cotacao.Bid)

}

func obterCotacao() common.Cotacao {
	// Definindo Context com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Criando requisição para o servidor
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	// Fazendo a requisição
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Decodificando o JSON
	return common.FromJsonToCotacao(b)

}
