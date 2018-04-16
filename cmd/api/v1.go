package api

import "github.com/erikvanbrakel/terraform-registry/cmd/registry"

type Disco struct {
	ModulesV1 string `json:"modules.v1"`
}

type ListModulesResponse struct {
	Meta    Meta              `json:"meta"`
	Modules []registry.Module `json:"modules"`
}

type ListVersionsResponse struct {
	Modules []registry.ModuleVersions `json:"modules"`
}

type Meta struct {
	Limit          int    `json:"limit"`
	PreviousOffset int    `json:"previous_offset,omitempty"`
	CurrentOffset  int    `json:"current_offset"`
	NextOffset     int    `json:"next_offset,omitempty"`
	NextUrl        string `json:"next_url,omitempty"`
	PreviousUrl    string `json:"prev_url,omitempty"`
}

type ApiError struct {
	Errors []string `json:"errors"`
}

func NewError(message string) ApiError {
	return ApiError {
		Errors: []string { message },
	}
}
