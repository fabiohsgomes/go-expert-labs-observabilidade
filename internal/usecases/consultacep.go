package usecases

import (
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/domain"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/clients"
)

type ConsultaCepUseCase struct {
	cepClient clients.CepClient
}

type DadosCep struct {
	Cep         string
	Logradouro  string
	Complemento string
	Bairro      string
	Localidade  string
	Uf          string
}

func NewConsultaCepUseCase(cepClient clients.CepClient) *ConsultaCepUseCase {
	return &ConsultaCepUseCase{
		cepClient: cepClient,
	}
}

func (u *ConsultaCepUseCase) ConsultaCep(cep *domain.Cep) (*DadosCep, error) {
	dadosCep, err := u.cepClient.ConsultaCep(cep.Codigo())
	if err != nil {
		return nil, err
	}

	return &DadosCep{
		Cep:         dadosCep.Cep,
		Logradouro:  dadosCep.Logradouro,
		Complemento: dadosCep.Complemento,
		Bairro:      dadosCep.Bairro,
		Localidade:  dadosCep.Localidade,
		Uf:          dadosCep.Uf,
	}, nil
}

func (u ConsultaCepUseCase) ValidateCep(cep string) bool {
	return len(cep) == 8
}
