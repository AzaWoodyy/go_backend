package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AzaWoodyy/go_backend/internal/models"
)

const (
	ddragonBaseURL = "https://ddragon.leagueoflegends.com"
)

type DDragonService struct {
	httpClient *http.Client
}

func NewDDragonService() *DDragonService {
	return &DDragonService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type DDragonChampion struct {
	ID    string   `json:"id"`
	Key   string   `json:"key"`
	Name  string   `json:"name"`
	Title string   `json:"title"`
	Blurb string   `json:"blurb"`
	Tags  []string `json:"tags"`
}

type DDragonResponse struct {
	Data map[string]DDragonChampion `json:"data"`
}

func (s *DDragonService) GetLatestVersion() (string, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/api/versions.json", ddragonBaseURL))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var versions []string
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return "", err
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no versions available")
	}

	return versions[0], nil
}

func (s *DDragonService) GetChampions(version string) ([]models.Champion, error) {
	url := fmt.Sprintf("%s/cdn/%s/data/en_US/champion.json", ddragonBaseURL, version)
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ddragonResp DDragonResponse
	if err := json.Unmarshal(body, &ddragonResp); err != nil {
		return nil, err
	}

	var champions []models.Champion
	for _, champ := range ddragonResp.Data {
		champions = append(champions, models.Champion{
			ID:    champ.ID,
			Key:   champ.Key,
			Name:  champ.Name,
			Title: champ.Title,
			Blurb: champ.Blurb,
		})
	}

	return champions, nil
}
