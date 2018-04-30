package models

type Module struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	Version   string `json:"version"`
}
