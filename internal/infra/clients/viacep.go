package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/helpers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type ViaCepClient struct {
	tracer trace.Tracer
	client http.Client
}

var viacepuri = "https://viacep.com.br/ws/"
var format = "/json/"

func NewViaCepClient(tracer trace.Tracer) *ViaCepClient {
	return &ViaCepClient{
		tracer: tracer,
		client: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
	}
}

func (c *ViaCepClient) ConsultaCep(cep string) (*DadosCepResponse, error) {
	_, span := otel.StartSpan(context.Background(), c.tracer, "ConsultaCep")
	defer span.End()
	
	if !helpers.ValidateZipCode(cep) {
		return nil, erros.ErrInvalidZipCode
	}

	otel.AddSpanEvent(span, "Iniciando consulta ViaCep", nil)

	dadosCep := &DadosCepResponse{}

	req, err := http.NewRequest("GET", viacepuri+cep+format, nil)
	if err != nil {
		otel.RecordSpanError(span, err)
		return dadosCep, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		otel.RecordSpanError(span, err)
		return dadosCep, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return dadosCep, fmt.Errorf("error fetching data: %s", resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &dadosCep)

	if len(dadosCep.Erro) > 0 {
		return dadosCep, erros.ErrZipCodeNotFound
	}

	dadosCep.Cep = helpers.NormalizeZipCode(dadosCep.Cep)

	otel.AddSpanEvent(span, "Cep localizado", nil)

	return dadosCep, nil
}
