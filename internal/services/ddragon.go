package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AzaWoodyy/go_backend/internal/models"
)

const (
	versionsURLTemplate = "https://ddragon.leagueoflegends.com/api/versions.json"
	championURLTemplate = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
	httpTimeout         = 10 * time.Second
)

type DDragonService struct {
	client *http.Client
}

func NewDDragonService() *DDragonService {
	return &DDragonService{
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

func (s *DDragonService) GetLatestVersion() (string, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, versionsURLTemplate, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch versions: status code %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read versions response body: %w", err)
	}

	var versions []string
	unmarshalErr := json.Unmarshal(bodyBytes, &versions)
	if unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal versions JSON: %w", unmarshalErr)
	}

	if len(versions) == 0 {
		return "", errors.New("no versions found in the response")
	}

	return versions[0], nil
}

func (s *DDragonService) GetChampions(version string) (*models.ChampionDataResponse, error) {
	ctx := context.Background()
	url := fmt.Sprintf(championURLTemplate, version)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch champions for version %s: %w", version, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch champions: status code %d, body: %s",
			resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read champions response body: %w", err)
	}

	var championData models.ChampionDataResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &championData)
	if unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal champions JSON: %w", unmarshalErr)
	}

	return &championData, nil
}
