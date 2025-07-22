package main

import (
	"log"
	"net/http"

	"github.com/rafabene/go-projects/cepweather-otel/serviceB/internal/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const port = ":8081"

func main() {
	mux := http.NewServeMux()
	log.Printf("Server started on port %s", port)
	h := otelhttp.WithRouteTag("/api/v1/weather", http.HandlerFunc(handlers.HandleCepWeather))
	mux.Handle("/api/v1/weather", h)
	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}

}
