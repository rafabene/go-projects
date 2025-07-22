package main

import (
	"log"
	"net/http"

	"github.com/rafabene/go-projects/cepweather-otel/serviceA/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const port = ":8080"

func main() {
	mux := http.NewServeMux()
	log.Printf("Server started on port %s", port)
	h := otelhttp.WithRouteTag("/api/v1/cep", http.HandlerFunc(handlers.HandleCepWeather))
	mux.Handle("/api/v1/cep", h)
	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}
