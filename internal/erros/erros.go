package erros

import (
	"errors"
)

var ErrInvalidZipCode = errors.New("invalid zipcode")
var ErrZipCodeNotFound = errors.New("can not find zipcode")
var ErrCityIsRequired = errors.New("city is required")
var ErrCityNotFound = errors.New("can not find city")
