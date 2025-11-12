package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"weather-api-lab/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWeatherUseCase struct {
	mock.Mock
}

func (m *MockWeatherUseCase) GetWeatherByZipcode(zipcode string) (*domain.Weather, error) {
	args := m.Called(zipcode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Weather), args.Error(1)
}

func TestWeatherHandlerGetWeather(t *testing.T) {
	tests := []struct {
		name           string
		zipcode        string
		mockWeather    *domain.Weather
		mockErr        error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "sucesso - CEP válido",
			zipcode: "26140040",
			mockWeather: &domain.Weather{
				TempC: 25.5,
				TempF: 77.9,
				TempK: 298.5,
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"temp_C":25.5,"temp_F":77.9,"temp_K":298.5}`,
		},
		{
			name:           "error - invalid zipcode",
			zipcode:        "123",
			mockWeather:    nil,
			mockErr:        domain.ErrInvalidZipcode,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"message":"invalid zipcode"}`,
		},
		{
			name:           "error - zipcode not found",
			zipcode:        "99999999",
			mockWeather:    nil,
			mockErr:        domain.ErrZipcodeNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"can not find zipcode"}`,
		},
		{
			name:           "error - weather not found",
			zipcode:        "26140040",
			mockWeather:    nil,
			mockErr:        domain.ErrWeatherNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"weather not found"}`,
		},
		{
			name:           "error - invalid location",
			zipcode:        "26140040",
			mockWeather:    nil,
			mockErr:        domain.ErrInvalidLocation,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid location"}`,
		},
		{
			name:           "error - internal server error",
			zipcode:        "26140040",
			mockWeather:    nil,
			mockErr:        assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar mock
			mockUseCase := new(MockWeatherUseCase)
			mockUseCase.On("GetWeatherByZipcode", tt.zipcode).Return(tt.mockWeather, tt.mockErr)

			// Criar handler
			handler := NewWeatherHandler(mockUseCase)
			router := handler.SetupRoutes()

			// Criar requisição
			req := httptest.NewRequest("GET", "/weather/"+tt.zipcode, nil)
			recorder := httptest.NewRecorder()

			// Executar requisição
			router.ServeHTTP(recorder, req)

			// Verificar resultado
			assert.Equal(t, tt.expectedStatus, recorder.Code)
			assert.JSONEq(t, tt.expectedBody, recorder.Body.String())
			assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

			// Verificar se o mock foi chamado
			mockUseCase.AssertExpectations(t)
		})
	}
}
