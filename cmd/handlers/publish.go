package handlers

import (
	"encoding/json"
	"github.com/erikvanbrakel/anthology/cmd/api"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"github.com/gorilla/mux"
	"net/http"
)

func PublishHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)

		namespace, name, provider, version := vars["namespace"], vars["name"], vars["provider"], vars["version"]

		err := r.PublishModule(namespace, name, provider, version, request.Body)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(api.NewError(err.Error()))
		}
	}
}
