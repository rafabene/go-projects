package common

import (
	"encoding/json"
	"log"
)

type Cotacao struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

// Estrutura que engloba a chave USDBRL
type Resposta struct {
	USDtoBRL Cotacao `json:"USDBRL"`
}

// Função que faz o unmarshal do json para a estrutura Cotacao
func FromJsonToCotacao(b []byte) Cotacao {
	var r *Resposta = &Resposta{}
	err := json.Unmarshal(b, r)
	if err != nil {
		log.Fatal(err)
	}
	return r.USDtoBRL
}
