package main

import (
	"log"
	"net/http"

	"github.com/rafabene/go-projects/cepweather/internal/handlers"
)

const port = ":8080"

func main() {
	log.Printf("Server started on port %s", port)
	http.HandleFunc("/api/v1/weather", handlers.HandleCepWeather)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
