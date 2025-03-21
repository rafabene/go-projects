package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rafabene/go-projects/multithreading/config"
)

var cfg = config.GetConfig()
var cep = cfg.Cep

type RequestData struct {
	Url   string
	Canal chan string
}

func main() {
	log.Println("Starting...")

	// Criando as estruturas de dados para as requisições 1
	rd1 := RequestData{
		Url:   cfg.GetUrl1(cep),
		Canal: make(chan string),
	}
	go doRequest(rd1)

	// Criando as estruturas de dados para as requisições 2
	rd2 := RequestData{
		Url:   cfg.GetUrl2(cep),
		Canal: make(chan string),
	}
	go doRequest(rd2)

	// Aguardando as respostas
	log.Println("Aguardando respostas...")
	select {
	case resp1 := <-rd1.Canal:
		imprimirReposta(rd1.Url, resp1)
	case resp2 := <-rd2.Canal:
		imprimirReposta(rd2.Url, resp2)
	case <-time.After(time.Duration(cfg.TimeoutMillis) * time.Millisecond):
		log.Println("Timeout!")
	}
	log.Println("Finalizando...")
}

func imprimirReposta(url string, resp string) {
	println("Resposta de ", url)
	println("Resposta: ", resp)
}

// Faz a requisição de acordo com o RequestData
// e preenchendo o canal com a resposta
func doRequest(rd RequestData) {
	resp, err := http.Get(rd.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Lendo o corpo da resposta
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Formatando a resposta
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")

	rd.Canal <- out.String()
}
