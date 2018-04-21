package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"github.com/erikvanbrakel/anthology/cmd/api"
)

func ListVersionsHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace, name, provider := params["namespace"], params["name"], params["provider"]

		modules, _ := r.ListVersions(namespace, name, provider)
		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(api.ListVersionsResponse{Modules: modules})
	}
}
