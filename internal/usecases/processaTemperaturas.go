package usecases

import (
	"log"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/domain"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/clients"
)

type DadosCepInput struct {
	Cep string `json:"cep"`
}

type DadosTemperaturasOutput struct {
	City       string `json:"city"`
	Celcius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}

type ProcessaTemperaturasService struct {
	client clients.CalculaTemperaturasClient
}

func NewProcessaTemperaturasService(client clients.CalculaTemperaturasClient) *ProcessaTemperaturasService {
	return &ProcessaTemperaturasService{
		client: client,
	}
}

func (s *ProcessaTemperaturasService) Execute(input DadosCepInput) (dados *DadosTemperaturasOutput, err error) {
	cep, err := domain.NewCep(input.Cep)
	if err != nil {
		return dados, err
	}

	response, err := s.client.CalculaTemperaturas(cep.Codigo())
	if err != nil {
		log.Println(err.Error())
		return dados, err
	}

	dados = &DadosTemperaturasOutput{
		City:       response.City,
		Celcius:    response.Celcius,
		Fahrenheit: response.Fahrenheit,
		Kelvin:     response.Kelvin,
	}

	return dados, err
}
