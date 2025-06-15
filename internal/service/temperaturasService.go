package service

import (
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/domain"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/usecases"
)

type TemperaturasService struct {
	consutlaCepUseCase         *usecases.ConsultaCepUseCase
	calculaTemperaturasUseCase *usecases.CalculaTemperaturasUseCase
}

func NewTemperaturasService(
	consutlaCepUseCase *usecases.ConsultaCepUseCase,
	calculaTemperaturasUseCase *usecases.CalculaTemperaturasUseCase,
) *TemperaturasService {
	return &TemperaturasService{
		consutlaCepUseCase:         consutlaCepUseCase,
		calculaTemperaturasUseCase: calculaTemperaturasUseCase,
	}
}

func (s *TemperaturasService) Processa(cep string) (*usecases.DadosTemperaturas, error) {
	cepDomain, err := domain.NewCep(cep)
	if err != nil {
		return nil, err
	}

	dadosCep, err := s.consutlaCepUseCase.ConsultaCep(cepDomain)
	if err != nil {
		return nil, err
	}

	localidadeDomain, err := domain.NewLocalidade(dadosCep.Localidade)
	if err != nil {
		return nil, err
	}

	dadosTemperaturas, err := s.calculaTemperaturasUseCase.Execute(localidadeDomain)
	if err != nil {
		return nil, err
	}

	return dadosTemperaturas, nil
}
