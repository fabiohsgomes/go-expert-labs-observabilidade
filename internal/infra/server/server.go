package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
)

type server struct {
	serviceName string
	port        int32
	mux         *http.ServeMux
}

func NewServer(serviceName string, port int32) *server {
	return &server{
		serviceName: serviceName,
		port:        port,
		mux:         http.NewServeMux(),
	}
}

func (s *server) Run(serverMuxCallBack func(serverMux *http.ServeMux)) {
	shutdown := otel.InitTracer(s.serviceName)
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("falha ao desligar o provedor de tracer: %v", err)
		}
	}()

	serverMuxCallBack(s.mux)
	log.Printf("Servidor escutando na porta :%d", s.port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux); err != nil {
		log.Println(err)
	}
}
