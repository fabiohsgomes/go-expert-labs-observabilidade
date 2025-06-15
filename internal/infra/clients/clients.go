package clients

import "fmt"

type CepClient interface {
	ConsultaCep(cep string) (*DadosCepResponse, error)
}

type DadosCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Erro        string `json:"erro,omitempty"`
}

type WeatherClient interface {
	ConsultaClima(cidade string) (*WeatherResponse, error)
}

type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int     `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Current struct {
	LastUpdatedEpoch int        `json:"last_updated_epoch"`
	LastUpdated      string     `json:"last_updated"`
	TempC            float64    `json:"temp_c"`
	TempF            float64    `json:"temp_f"`
	IsDay            int        `json:"is_day"`
	Condition        Condition  `json:"condition"`
	WindMph          float64    `json:"wind_mph"`
	WindKph          float64    `json:"wind_kph"`
	WindDegree       int        `json:"wind_degree"`
	WindDir          string     `json:"wind_dir"`
	PressureMb       int        `json:"pressure_mb"`
	PressureIn       float64    `json:"pressure_in"`
	PrecipMm         int        `json:"precip_mm"`
	PrecipIn         int        `json:"precip_in"`
	Humidity         int        `json:"humidity"`
	Cloud            int        `json:"cloud"`
	FeelslikeC       int        `json:"feelslike_c"`
	FeelslikeF       float64    `json:"feelslike_f"`
	VisKm            int        `json:"vis_km"`
	VisMiles         int        `json:"vis_miles"`
	Uv               int        `json:"uv"`
	GustMph          float64    `json:"gust_mph"`
	GustKph          float64    `json:"gust_kph"`
	AirQuality       Airquality `json:"air_quality"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Airquality struct {
	Co           float64 `json:"co"`
	No2          float64 `json:"no2"`
	O3           float64 `json:"o3"`
	So2          int     `json:"so2"`
	Pm25         float64 `json:"pm2_5"`
	Pm10         int     `json:"pm10"`
	UsEpaIndex   int     `json:"us-epa-index"`
	GbDefraIndex int     `json:"gb-defra-index"`
}

type WeatherErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewWeatherErrorResponse(code int, message string) WeatherErrorResponse {
	return WeatherErrorResponse{
		Code:    code,
		Message: message,
	}
}

func (e WeatherErrorResponse) ErrorCode() int {
	return e.Code
}

func (e WeatherErrorResponse) ErrorMessage() string {
	return e.Message
}

func (e WeatherErrorResponse) Error() string {
	return fmt.Sprintf("%d :: %s", e.Code, e.Message)
}

type CalculaTemperaturasClient interface {
	CalculaTemperaturas(cep string) (*TemperaturasResponse, error)
}

type TemperaturasResponse struct {
	City       string `json:"city"`
	Celcius    string `json:"temp_C"`
	Fahrenheit string `json:"temp_F"`
	Kelvin     string `json:"temp_K"`
}
