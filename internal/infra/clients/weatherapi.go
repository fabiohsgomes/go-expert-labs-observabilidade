package clients

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/config"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

type WeatherApiClient struct {
	tracer trace.Tracer
	client http.Client
	key    string
}

var weatherapiuri = "https://api.weatherapi.com/v1/current.json"

func NewWeatherApiClient(tracer trace.Tracer) *WeatherApiClient {
	cfg := config.Get()

	return &WeatherApiClient{
		tracer: tracer,
		client: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   10 * time.Second,
		},
		key: cfg.GetWeatherApiKey(),
	}
}

func (c *WeatherApiClient) ConsultaClima(cidade string) (*WeatherResponse, error) {
	_, span := otel.StartSpan(context.Background(), c.tracer, "ConsultaClima")
	defer span.End()

	otel.AddSpanEvent(span, "Iniciando a consulta WeatherAPI", nil)

	weatherResponse := &WeatherResponse{}
	weatherErrorResponse := WeatherErrorResponse{}

	req, err := http.NewRequest("GET", weatherapiuri, nil)
	if err != nil {
		otel.RecordSpanError(span, err)
		return weatherResponse, err
	}

	req.Header.Set("Accept", "application/json")
	url := req.URL
	q := url.Query()
	q.Set("q", cidade)
	q.Set("lang", "pt")
	q.Set("key", c.key)
	url.RawQuery = q.Encode()

	req.URL = url

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		otel.RecordSpanError(span, err)
		return weatherResponse, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		json.Unmarshal(body, &weatherErrorResponse)

		if weatherErrorResponse.ErrorCode() == 1006 {
			return weatherResponse, erros.ErrCityNotFound
		}

		return weatherResponse, weatherErrorResponse
	}

	json.Unmarshal(body, &weatherResponse)

	otel.AddSpanEvent(span, "Temperaturas obtidas com sucesso", nil)

	return weatherResponse, nil
}
