package domain

import (
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/helpers"
)

type Cep struct {
	codigo string
}

func NewCep(codigo string) (*Cep, error) {
	cep := &Cep{
		codigo: helpers.NormalizeZipCode(codigo),
	}

	if !helpers.ValidateZipCode(cep.codigo) {
		return nil, erros.ErrInvalidZipCode
	}

	return cep, nil
}

func (c *Cep) Codigo() string {
	return c.codigo
}
