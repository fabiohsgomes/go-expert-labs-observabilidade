package domain

import "github.com/fabiohsgomes/go-expert-labs-deploy/internal/erros"

type Localidade struct {
	name string
}

func NewLocalidade(name string) (*Localidade, error) {
	localidade := &Localidade{
		name: name,
	}

	if !localidade.valida() {
		return nil, erros.ErrCityIsRequired
	}

	return localidade, nil
}

func (l *Localidade) Name() string {
	return l.name
}

func (l Localidade) valida() bool {
	return len(l.name) > 0
}
