package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"weather-api-lab/internal/domain"
	"weather-api-lab/internal/dto"
)

type ViaCEPClient interface {
	GetLocationByZipcode(zipcode string) (*domain.Location, error)
}

type viacepClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewViaCEPClient(baseURL string) ViaCEPClient {
	return &viacepClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *viacepClient) GetLocationByZipcode(zipcode string) (*domain.Location, error) {
	if err := domain.ValidateZipcode(zipcode); err != nil {
		return nil, err
	}

	formattedZipcode := domain.FormatZipcode(zipcode)
	url := fmt.Sprintf("%s/%s/json/", c.baseURL, formattedZipcode)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error querying ViaCEP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading ViaCEP response: %w", err)
	}

	var viacepResp dto.ViaCEPResponse
	if err := json.Unmarshal(body, &viacepResp); err != nil {
		return nil, fmt.Errorf("error parsing ViaCEP response: %w", err)
	}

	if viacepResp.Erro == "true" || viacepResp.Localidade == "" {
		return nil, domain.ErrZipcodeNotFound
	}

	return &domain.Location{
		City:  viacepResp.Localidade,
		State: viacepResp.UF,
	}, nil
}
