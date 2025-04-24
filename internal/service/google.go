package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

var (
	autocompleteMinLength            = 3
	ErrGoogleAutocompleteQueryLength = fmt.Errorf("query must be at least %d characters long", autocompleteMinLength)
	ErrGoogleServiceUnavailable      = fmt.Errorf("google service is unavailable")
)

type GoogleService interface {
	Autocomplete(ctx context.Context, query string) ([]dto.GooglePredictionDTO, error)
}

type googleService struct {
	cfg *config.Google
}

func NewGoogleService(cfg *config.Google) *googleService {
	return &googleService{
		cfg: cfg,
	}
}

func (s *googleService) Autocomplete(ctx context.Context, query string) ([]dto.GooglePredictionDTO, error) {
	if len(query) < 3 {
		return nil, ErrGoogleAutocompleteQueryLength
	}

	const apiURL = "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Set("input", query)
	params.Set("key", s.cfg.APIKey)
	reqURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("google autocomplete API returned status: %s", resp.Status)
		return nil, ErrGoogleServiceUnavailable
	}

	var result struct {
		Predictions []struct {
			PlaceID              string `json:"place_id"`
			Description          string `json:"description"`
			StructuredFormatting struct {
				MainText      string `json:"main_text"`
				SecondaryText string `json:"secondary_text"`
			} `json:"structured_formatting"`
		} `json:"predictions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding Google API response: %v", err)
		return nil, err
	}

	var suggestions []dto.GooglePredictionDTO
	for _, prediction := range result.Predictions {
		suggestions = append(suggestions, dto.GooglePredictionDTO{
			PlaceID:     prediction.PlaceID,
			Description: prediction.Description,
		})
	}

	if len(suggestions) == 0 {
		return []dto.GooglePredictionDTO{}, nil
	}

	return suggestions, nil
}
