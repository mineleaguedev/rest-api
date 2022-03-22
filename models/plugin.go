package models

type Plugin struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type PluginsResponse struct {
	Success bool     `json:"success"`
	Plugins []Plugin `json:"plugins"`
}

type PluginResponse struct {
	Success  bool     `json:"success"`
	Versions []string `json:"versions"`
}
