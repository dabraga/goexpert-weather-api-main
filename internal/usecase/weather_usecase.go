package usecase

import (
	"weather-api-lab/internal/domain"
	"weather-api-lab/internal/repository"
)

type WeatherUseCase interface {
	GetWeatherByZipcode(zipcode string) (*domain.Weather, error)
}

type weatherUseCase struct {
	viacepClient  repository.ViaCEPClient
	weatherClient repository.WeatherClient
}

func NewWeatherUseCase(viacepClient repository.ViaCEPClient, weatherClient repository.WeatherClient) WeatherUseCase {
	return &weatherUseCase{
		viacepClient:  viacepClient,
		weatherClient: weatherClient,
	}
}

func (u *weatherUseCase) GetWeatherByZipcode(zipcode string) (*domain.Weather, error) {
	// 1. Buscar localização pelo CEP
	location, err := u.viacepClient.GetLocationByZipcode(zipcode)
	if err != nil {
		return nil, err
	}

	// 2. Buscar temperatura pela localização
	tempCelsius, err := u.weatherClient.GetTemperatureByLocation(location)
	if err != nil {
		return nil, err
	}

	// 3. Criar objeto Weather com conversões
	weather := domain.NewWeather(tempCelsius)

	return &weather, nil
}
