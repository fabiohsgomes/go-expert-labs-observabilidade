package usecases

import (
	"fmt"
	"testing"

	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/domain"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/helpers"
	"github.com/fabiohsgomes/go-expert-labs-deploy/internal/infra/clients"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CalculaTemperaturasTestSuite struct {
	suite.Suite
}

type WeatherApiClientMock struct {
	mock.Mock
}

func (m *WeatherApiClientMock) ConsultaClima(cidade string) (*clients.WeatherResponse, error) {
	args := m.Called(cidade)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*clients.WeatherResponse), args.Error(1)
}

func TestCalculaTemperaturasSuite(t *testing.T) {
	suite.Run(t, new(CalculaTemperaturasTestSuite))
}

func (s *CalculaTemperaturasTestSuite) TestCalculaTemperaturas() {
	//Arrange
	weatherApiClientMock := new(WeatherApiClientMock)
	calculaTemperaturasUseCase := NewCalculaTemperaturasUseCase(weatherApiClientMock)

	cidade, _ := domain.NewLocalidade("São Paulo")
	expectedResponse := &clients.WeatherResponse{
		Current: clients.Current{
			TempC: 25.0,
		},
	}

	expectedCelsius := fmt.Sprintf("%.1f", expectedResponse.Current.TempC)
	expectedFahrenheit := fmt.Sprintf("%.1f", helpers.CelsiusToFahrenheit(expectedResponse.Current.TempC))
	expectedKelvin := fmt.Sprintf("%.1f", helpers.CelsiusToKelvin(expectedResponse.Current.TempC))

	weatherApiClientMock.On("ConsultaClima", cidade.Name()).Return(expectedResponse, nil)

	//Act
	dadosTemperaturas, err := calculaTemperaturasUseCase.Execute(cidade)

	//Assert
	s.NoError(err)
	s.Equal(expectedCelsius, dadosTemperaturas.Celcius)
	s.Equal(expectedFahrenheit, dadosTemperaturas.Fahrenheit) // 25°C to °F
	s.Equal(expectedKelvin, dadosTemperaturas.Kelvin)         // 25°C to K

	weatherApiClientMock.AssertExpectations(s.T())
}
