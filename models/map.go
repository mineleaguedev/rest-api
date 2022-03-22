package models

type MiniGames struct {
	Name    string   `json:"name"`
	Formats []Format `json:"formats"`
}

type Format struct {
	Format string `json:"format"`
	Maps   []Map  `json:"maps"`
}

type Map struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type MapsResponse struct {
	Success   bool        `json:"success"`
	MiniGames []MiniGames `json:"minigames"`
}

type MiniGameMapsResponse struct {
	Success bool     `json:"success"`
	Formats []Format `json:"formats"`
}

type MiniGameFormatMapsResponse struct {
	Success bool  `json:"success"`
	Maps    []Map `json:"maps"`
}
