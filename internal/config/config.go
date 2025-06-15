package config

//import viper
import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type configApp struct {
	AmbientePublicacao string `mapstructure:"AMBIENTE_PUBLICACAO"`
	WeatherApiKey      string `mapstructure:"WEATHER_API_KEY"`
}

var config *configApp

func LoadConfig(path string) {
	if _, err := os.Stat(filepath.Join(path, ".env")); err == nil {
		viper.SetConfigName(".env")
		viper.SetConfigType("env")
		viper.AddConfigPath(path)
		viper.AutomaticEnv()

		err := viper.ReadInConfig()
		if err != nil {
			panic(err)
		}

		err = viper.Unmarshal(&config)
		if err != nil {
			panic(err)
		}

		return
	}

	var ambientepublicacao string
	if value, exists := os.LookupEnv("AMBIENTE_PUBLICACAO"); !exists {
		panic("variabel de ambiente AMBIENTE_PUBLICACAO, não definida")
	} else {
		ambientepublicacao = value
	}

	var weatherapikey string
	if value, exists := os.LookupEnv("WEATHER_API_KEY"); !exists {
		panic("variavel de ambiente WEATHER_API_KEY, não definida")
	} else {
		weatherapikey = value
	}
	config = &configApp{}
	config.AmbientePublicacao = ambientepublicacao
	config.WeatherApiKey = weatherapikey

	return
}

func Get() *configApp {
	if config == nil{
		LoadConfig(".")
	}

	return config
}

func (c *configApp) GetAmbientePublicacao() string {
	return c.AmbientePublicacao
}

func (c *configApp) GetWeatherApiKey() string {
	return c.WeatherApiKey
}
