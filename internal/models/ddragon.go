package models

type ChampionInfo struct {
	Version string `json:"version"`
	Key     string `json:"key"`
	Name    string `json:"name"`
	Title   string `json:"title"`
}

type ChampionDataResponse struct {
	Type    string                  `json:"type"`
	Format  string                  `json:"format"`
	Version string                  `json:"version"`
	Data    map[string]ChampionInfo `json:"data"`
}
