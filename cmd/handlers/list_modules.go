package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
	"github.com/erikvanbrakel/terraform-registry/cmd/api"
)

func ListModulesHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func (writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace := params["namespace"]

		modules, _ := r.ListModules(namespace, "", "", 0, 99999)

		meta := api.Meta{
			CurrentOffset: 0,
			Limit:         9999,
			NextOffset:    2,
			NextUrl:       "",
		}

		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(api.ListModulesResponse{Meta: meta, Modules: modules})
	}
}
