package otel

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer configura o provedor de tracer OpenTelemetry para enviar traces via OTLP.
func InitTracer(serviceName string) func(context.Context) error {
	collectorEndpoint := os.Getenv("OTEL_COLLECTOR_ENDPOINT")
	if collectorEndpoint == "" {
		collectorEndpoint = "localhost:4317" // Endereço padrão do gRPC do OTel Collector
	}

	// NOVO: Crie um contexto com tempo limite para a conexão gRPC.
	// Isso substitui grpc.WithBlock() e grpc.WithTimeout().
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Certifique-se de cancelar o contexto quando ele não for mais necessário

	// Configura o cliente gRPC para o exportador OTLP
	// As opções grpc.WithBlock() e grpc.WithTimeout() foram removidas.
	conn, err := grpc.NewClient(
		collectorEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Usar credenciais seguras em produção
		// Não precisamos mais de WithBlock ou WithTimeout aqui, pois o contexto faz isso.
		// grpc.WithBlock(), // DEPRECATED
		// grpc.WithTimeout(5*time.Second), // DEPRECATED
	)
	if err != nil {
		log.Fatalf("falha ao criar conexão gRPC para o OTLP Collector: %v", err)
	}

	// Cria o exportador OTLP gRPC
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("falha ao criar o exportador OTLP: %v", err)
	}

	// Cria o provedor de tracer
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("environment", "development"),
			attribute.String("application", "zipcode-weather-app"),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// A função de shutdown precisa garantir que o provedor de tracer seja desligado e a conexão fechada.
	return func(ctx context.Context) error {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Erro ao desligar o provedor de tracer: %v", err)
			// Não retorne imediatamente, tente fechar a conexão também.
		}
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Printf("Erro ao fechar conexão gRPC: %v", err)
				return err // Retorne o erro de fechamento de conexão se o provedor já desligou sem erro
			}
		}
		return nil
	}
}

// GetTracer retorna um tracer para o serviço especificado.
func GetTracer(serviceName string) trace.Tracer {
	return otel.Tracer(serviceName)
}

// StartSpan inicia um novo span e o retorna juntamente com o contexto.
func StartSpan(ctx context.Context, tracer trace.Tracer, spanName string) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, spanName)
	return ctx, span
}

// AddSpanEvent adiciona um evento a um span.
func AddSpanEvent(span trace.Span, name string, attributes map[string]interface{}) {
	if span == nil {
		return
	}
	// Converte o mapa de interface{} para []attribute.KeyValue
	otelAttributes := make([]attribute.KeyValue, 0, len(attributes))
	for k, v := range attributes {
		switch val := v.(type) {
		case string:
			otelAttributes = append(otelAttributes, attribute.String(k, val))
		case int:
			otelAttributes = append(otelAttributes, attribute.Int(k, val))
		case float64:
			otelAttributes = append(otelAttributes, attribute.Float64(k, val))
		case bool:
			otelAttributes = append(otelAttributes, attribute.Bool(k, val))
		default:
			otelAttributes = append(otelAttributes, attribute.String(k, fmt.Sprintf("%v", val)))
		}
	}
	span.AddEvent(name, trace.WithAttributes(otelAttributes...))
}

// RecordSpanError registra um erro em um span.
func RecordSpanError(span trace.Span, err error) {
	if span == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
