package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/memsbdm/restaurant-api/config"
	"github.com/memsbdm/restaurant-api/internal/dto"
)

var (
	autocompleteMinLength            = 3
	ErrGoogleAutocompleteQueryLength = fmt.Errorf("query must be at least %d characters long", autocompleteMinLength)
	ErrGoogleServiceUnavailable      = fmt.Errorf("google service is unavailable")
	ErrGoogleInvalidPlaceID          = fmt.Errorf("invalid place ID")
)

type GoogleService interface {
	Autocomplete(ctx context.Context, query string) ([]*dto.GooglePrediction, error)
	GetDetails(ctx context.Context, placeID string) (*dto.CreateRestaurant, error)
}

type googleService struct {
	cfg    *config.Google
	client *http.Client
}

func NewGoogleService(cfg *config.Google) *googleService {
	return &googleService{
		cfg: cfg,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *googleService) Autocomplete(ctx context.Context, query string) ([]*dto.GooglePrediction, error) {
	if len(query) < 3 {
		return nil, ErrGoogleAutocompleteQueryLength
	}

	const apiURL = "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Set("input", query)
	params.Set("key", s.cfg.APIKey)
	reqURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	resp, err := s.client.Get(reqURL)
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

	predictions := make([]*dto.GooglePrediction, len(result.Predictions))
	for i := range result.Predictions {
		predictions[i] = &dto.GooglePrediction{
			PlaceID:     result.Predictions[i].PlaceID,
			Description: result.Predictions[i].Description,
		}
	}

	if len(predictions) == 0 {
		return []*dto.GooglePrediction{}, nil
	}

	return predictions, nil
}

func (s *googleService) GetDetails(ctx context.Context, placeID string) (*dto.CreateRestaurant, error) {
	const apiURL = "https://places.googleapis.com/v1/places/"
	params := url.Values{}
	params.Set("fields", "displayName,formattedAddress,location,internationalPhoneNumber")
	params.Set("key", s.cfg.APIKey)
	reqURL := fmt.Sprintf("%s%s?%s", apiURL, placeID, params.Encode())

	resp, err := s.client.Get(reqURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return nil, ErrGoogleInvalidPlaceID
		}
		log.Printf("google place details API returned status: %s", resp.Status)
		return nil, ErrGoogleServiceUnavailable
	}

	var result struct {
		InternationalPhoneNumber *string `json:"internationalPhoneNumber"`
		FormattedAddress         string  `json:"formattedAddress"`
		Location                 *struct {
			Lat *float64 `json:"lat"`
			Lng *float64 `json:"lng"`
		} `json:"location"`
		DisplayName struct {
			Text        string `json:"text"`
			LangageCode string `json:"languageCode"`
		} `json:"displayName"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding Google API response: %v", err)
		return nil, err
	}

	restaurant := &dto.CreateRestaurant{
		Name:    result.DisplayName.Text,
		Alias:   result.DisplayName.Text,
		Address: result.FormattedAddress,
		PlaceID: placeID,
	}

	if result.Location != nil {
		restaurant.Lat = result.Location.Lat
		restaurant.Lng = result.Location.Lng
	}

	formattedPhone := strings.ReplaceAll(*result.InternationalPhoneNumber, " ", "")
	restaurant.Phone = &formattedPhone

	return restaurant, nil
}
