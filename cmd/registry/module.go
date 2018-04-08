package registry

type Module struct {
	ID          string `json:"id"`
	Owner       string `json:"owner"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Provider    string `json:"provider"`
	Description string `json:"description"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
	Downloads   int    `json:"downloads"`
	Verified    bool   `json:"verified"`
}

type ModuleVersions struct {
	Source   string   `json:"source"`
	Versions []Module `json:"versions"`
}
