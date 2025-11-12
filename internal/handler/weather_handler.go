package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"weather-api-lab/internal/domain"
	"weather-api-lab/internal/dto"
	"weather-api-lab/internal/usecase"
)

type WeatherHandler struct {
	weatherUseCase usecase.WeatherUseCase
}

func NewWeatherHandler(weatherUseCase usecase.WeatherUseCase) *WeatherHandler {
	return &WeatherHandler{
		weatherUseCase: weatherUseCase,
	}
}

func (h *WeatherHandler) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Get("/weather/{zipcode}", h.GetWeather)

	return r
}

// GetWeather busca o clima pelo CEP
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	zipcode := chi.URLParam(r, "zipcode")

	// Buscar clima
	weather, err := h.weatherUseCase.GetWeatherByZipcode(zipcode)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Retornar sucesso
	h.writeJSONResponse(w, http.StatusOK, weather)
}

// handleError trata erros e retorna resposta apropriada
func (h *WeatherHandler) handleError(w http.ResponseWriter, err error) {
	log.Printf("Error processing request: %v", err)

	switch err {
	case domain.ErrInvalidZipcode:
		h.writeErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
	case domain.ErrZipcodeNotFound:
		h.writeErrorResponse(w, http.StatusNotFound, "can not find zipcode")
	case domain.ErrWeatherNotFound:
		h.writeErrorResponse(w, http.StatusNotFound, "weather not found")
	case domain.ErrInvalidLocation:
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid location")
	default:
		h.writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *WeatherHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func (h *WeatherHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	errorResp := dto.ErrorResponse{
		Message: message,
	}
	h.writeJSONResponse(w, statusCode, errorResp)
}
