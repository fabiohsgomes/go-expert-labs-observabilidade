package usecases

import (
	"fmt"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/domain"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/helpers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/clients"
)

type CalculaTemperaturasUseCase struct {
	weatherapiClient clients.WeatherClient
}

func NewCalculaTemperaturasUseCase(weatherapiClient clients.WeatherClient) *CalculaTemperaturasUseCase {
	return &CalculaTemperaturasUseCase{
		weatherapiClient: weatherapiClient,
	}
}

type DadosTemperaturas struct {
	City       string `json:"city"`
	Celcius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

func (u *CalculaTemperaturasUseCase) Execute(localidade *domain.Localidade) (*DadosTemperaturas, error) {
	weatherResponse, err := u.weatherapiClient.ConsultaClima(localidade.Name())
	if err != nil {
		return nil, err
	}

	dadosTemperaturas := u.processaTemperaturas(weatherResponse)
	dadosTemperaturas.City = localidade.Name()

	return dadosTemperaturas, nil
}

func (u *CalculaTemperaturasUseCase) processaTemperaturas(weatherResponse *clients.WeatherResponse) *DadosTemperaturas {
	return &DadosTemperaturas{
		Celcius:    fmt.Sprintf("%.1f", weatherResponse.Current.TempC),
		Fahrenheit: fmt.Sprintf("%.1f", helpers.CelsiusToFahrenheit(weatherResponse.Current.TempC)),
		Kelvin:     fmt.Sprintf("%.1f", helpers.CelsiusToKelvin(weatherResponse.Current.TempC)),
	}
}
