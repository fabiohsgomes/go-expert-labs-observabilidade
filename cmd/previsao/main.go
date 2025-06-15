package main

import (
	"net/http"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/config"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/handlers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/server"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var serviceName string = "Serviço B"

func main() {
	config.LoadConfig(".")

	server := server.NewServer(serviceName, 3001)

	server.Run(func(mux *http.ServeMux) {
		tracer := otel.GetTracer(serviceName)

		mux.Handle("GET /cidades/{cep}/temperaturas", otelhttp.NewHandler(http.HandlerFunc(handlers.ProcessaTemperaturasHandler(tracer)), "handlerServic-B"))
	})
}
