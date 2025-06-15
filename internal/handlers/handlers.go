package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/helpers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/clients"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/service"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/usecases"
	"github.com/fabiohsgomes/go-expert-labs-deploy/pkg/observabilidade/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func CapturaTemperaturasHandler(tracer trace.Tracer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.StartSpan(r.Context(), tracer, "CapturaTemperaturasHandler")
		otel.AddSpanEvent(span, "Recebendo requisição de CEP", nil)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dadosInput := usecases.DadosCepInput{}

		json.Unmarshal(body, &dadosInput)

		if !helpers.ValidateZipCode(dadosInput.Cep) {
			span.SetStatus(codes.Error, erros.ErrInvalidZipCode.Error())
			http.Error(w, erros.ErrInvalidZipCode.Error(), http.StatusUnprocessableEntity)
			return
		}

		otel.AddSpanEvent(span, "Cep validado, encaminhando para o Serviço B", map[string]interface{}{"cep": dadosInput.Cep})
		calculaTemperaturasClient := clients.NewCalculaTemperaturasClient(tracer)
		service := usecases.NewProcessaTemperaturasService(calculaTemperaturasClient)

		dados, err := service.Execute(dadosInput)
		if err != nil {
			if errors.Is(err, erros.ErrInvalidZipCode) {
				span.SetStatus(codes.Error, err.Error())
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}

			if errors.Is(err, erros.ErrZipCodeNotFound) {
				span.SetStatus(codes.Error, err.Error())
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			otel.RecordSpanError(span, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Header.Set(w.Header(), "Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dados); err != nil {
			otel.RecordSpanError(span, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error encoding response for CEP %s: %v", dadosInput.Cep, err)
			return
		}

		span.SetStatus(codes.Ok, "Requisição processada com sucesso")
	}
}

func ProcessaTemperaturasHandler(tracer trace.Tracer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, span := otel.StartSpan(r.Context(), tracer, "ProcessaTemperaturasHandler")

		otel.AddSpanEvent(span, "Recebendo requisição de processamendo do clima", nil)

		cepPathValue := r.PathValue("cep")
		viaCepClient := clients.NewViaCepClient(tracer)
		cepUseCase := usecases.NewConsultaCepUseCase(viaCepClient)

		weatherApiClient := clients.NewWeatherApiClient(tracer)
		calculaTemperaturasUseCase := usecases.NewCalculaTemperaturasUseCase(weatherApiClient)

		temperaturasService := service.NewTemperaturasService(cepUseCase, calculaTemperaturasUseCase)
		dadosTemperaturas, err := temperaturasService.Processa(cepPathValue)
		if err != nil {
			if errors.Is(err, erros.ErrInvalidZipCode) || errors.Is(err, erros.ErrCityIsRequired) {
				span.SetStatus(codes.Error, err.Error())
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}
			if errors.Is(err, erros.ErrZipCodeNotFound) || errors.Is(err, erros.ErrCityNotFound) {
				span.SetStatus(codes.Error, err.Error())
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			span.SetStatus(codes.Error, err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error processing temperatures for CEP %s: %v", cepPathValue, err)
			return
		}

		span.SetStatus(codes.Ok, "Calculo das temperaturas realizada com sucesso")

		http.Header.Set(w.Header(), "Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(dadosTemperaturas); err != nil {
			otel.RecordSpanError(span, err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Printf("Error encoding response for CEP %s: %v", cepPathValue, err)
			return
		}

		log.Printf("Successfully processed temperatures for CEP %s", cepPathValue)
	}
}
