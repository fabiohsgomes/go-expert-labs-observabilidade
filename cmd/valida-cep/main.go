package main

import (
	"net/http"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/handlers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/server"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var serviceName = "Serviço A"

func main() {
	server := server.NewServer(serviceName, 3000)

	server.Run(func(mux *http.ServeMux) {
		tracer := otel.GetTracer(serviceName)

		mux.Handle("POST /temperaturas", otelhttp.NewHandler(http.HandlerFunc(handlers.CapturaTemperaturasHandler(tracer)), "handlerServic-A"))
	})
}
