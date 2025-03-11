package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rafabene/go-projects/client-server-api/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const port = ":8080"
const path = "/cotacao"

var db *gorm.DB

func main() {

	db, _ = gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})

	// Migrate the schema
	db.AutoMigrate(&common.Cotacao{})

	// Iniciar o servidor
	startServer()
}

// Função que inicia o servidor
func startServer() {
	defer http.ListenAndServe(port, nil)
	http.HandleFunc(path, obterCotacao)
	log.Printf("Servidor rodando na porta %s", port)
}

func obterCotacao(w http.ResponseWriter, r *http.Request) {
	//Declarando contexto com timeout
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	//Fazer a requisição para a API de cotação
	log.Println("Obtendo cotação")
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if (err != nil) || (resp.StatusCode != http.StatusOK) {
		var errMsg string = "Erro fazendo a requisição."
		http.Error(w, errMsg+err.Error(), http.StatusInternalServerError)
		log.Println(errMsg, err.Error())
		return
	}
	defer resp.Body.Close()

	//Ler o corpo da resposta
	b, err := io.ReadAll(resp.Body)
	if (err != nil) || (len(b) == 0) {
		var errMsg string = "Erro lendo o corpo da resposta."
		http.Error(w, errMsg+err.Error(), http.StatusInternalServerError)
		log.Println(errMsg, err.Error())
		return
	}
	log.Printf("Cotação obtida com sucesso: %s", b)

	//Fazer o unmarshal do json para a estrutura Cotacao
	cotacao := common.FromJsonToCotacao(b)

	//Mostrar cotação na resposta
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)

	//Registrar cotação no banco de dados
	registrarCotacao(cotacao)
}

// Função que simula o registro da cotação no banco de dados
func registrarCotacao(cotacao common.Cotacao) {
	log.Printf("Registrando Cotação no banco de dados.")

	// Contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result := db.WithContext(ctx).Debug().Create(&cotacao)

	println(result.Error)
}
