package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"weather-api-lab/internal/handler"
	"weather-api-lab/internal/repository"
	"weather-api-lab/internal/usecase"
)

func main() {
	config := setupConfig()

	viacepClient := repository.NewViaCEPClient(config.GetString("viacep_base_url"))
	weatherClient := repository.NewWeatherClient(config.GetString("weather_api_base_url"), config.GetString("weather_api_key"))
	weatherUseCase := usecase.NewWeatherUseCase(viacepClient, weatherClient)
	weatherHandler := handler.NewWeatherHandler(weatherUseCase)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GetInt("port")),
		Handler:      weatherHandler.SetupRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("Server running on port %d", config.GetInt("port"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error stopping server: %v", err)
	}

	log.Println("Server stopped successfully")
}

func setupConfig() *viper.Viper {
	v := viper.New()

	v.SetDefault("port", 8080)
	v.SetDefault("weather_api_key", "7799569b0a824f369b504330251211")
	v.SetDefault("weather_api_base_url", "https://api.weatherapi.com/v1")
	v.SetDefault("viacep_base_url", "https://viacep.com.br/ws")

	v.AutomaticEnv()

	v.SetConfigFile(".env")
	if err := v.ReadInConfig(); err == nil {
		log.Println(".env file loaded")
	}

	if v.GetString("weather_api_key") == "" {
		log.Fatal("WEATHER_API_KEY is required. Configure in .env or as environment variable")
	}

	return v
}
