package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type CalculaTemperaturasClientService struct {
	tracer trace.Tracer
	uri    string
	client http.Client
}

func NewCalculaTemperaturasClient(tracer trace.Tracer) *CalculaTemperaturasClientService {
	return &CalculaTemperaturasClientService{
		tracer: tracer,
		uri:    "http://service-b:3001/",
		client: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
	}
}

func (c *CalculaTemperaturasClientService) CalculaTemperaturas(cep string) (response *TemperaturasResponse, err error) {
	_, span := otel.StartSpan(context.Background(), c.tracer, "CalculaTemperaturas")
	defer span.End()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%scidades/%s/temperaturas", c.uri, cep), nil)
	if err != nil {
		otel.RecordSpanError(span, err)
		return response, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		otel.RecordSpanError(span, err)
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnprocessableEntity {
			return response, erros.ErrInvalidZipCode
		}

		if resp.StatusCode == http.StatusNotFound {
			return response, erros.ErrZipCodeNotFound
		}

		return response, fmt.Errorf("error fetching data: %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &response)

	return response, err
}
