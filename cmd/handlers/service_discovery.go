package handlers

import (
	"encoding/json"
	"github.com/erikvanbrakel/terraform-registry/cmd/api"
	"net/http"
)

func ServiceDiscoveryHandler() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(api.Disco{ModulesV1: "/v1/modules/"})
	}
}
