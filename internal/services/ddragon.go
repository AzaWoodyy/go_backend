package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/AzaWoodyy/go_backend/internal/models"

	_ "log"
)

const (
	versionsURLTemplate = "https://ddragon.leagueoflegends.com/api/versions.json"
	championURLTemplate = "https://ddragon.leagueoflegends.com/cdn/%s/data/en_US/champion.json"
)

type DDragonService struct {
	client *http.Client
}

func NewDDragonService() *DDragonService {
	return &DDragonService{
		client: &http.Client{},
	}
}

func (s *DDragonService) GetLatestVersion() (string, error) {
	resp, err := s.client.Get(versionsURLTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to fetch versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch versions: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read versions response body: %w", err)
	}

	var versions []string
	if err := json.Unmarshal(bodyBytes, &versions); err != nil {
		return "", fmt.Errorf("failed to unmarshal versions JSON: %w", err)
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions found in the response")
	}

	return versions[0], nil
}

func (s *DDragonService) GetChampions(version string) (*models.ChampionDataResponse, error) {
	url := fmt.Sprintf(championURLTemplate, version)
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch champions for version %s: %w", version, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch champions: status code %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read champions response body: %w", err)
	}

	var championData models.ChampionDataResponse
	if err := json.Unmarshal(bodyBytes, &championData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal champions JSON: %w", err)
	}

	return &championData, nil
}
