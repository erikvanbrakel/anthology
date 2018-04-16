package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
	"github.com/erikvanbrakel/terraform-registry/cmd/api"
	"strconv"
)

func ListModulesHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func (writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		query := request.URL.Query()

		offset,_ := strconv.Atoi(query.Get("offset"))
		limit,err := strconv.Atoi(query.Get("limit"))

		if err != nil {
			limit = 10
		}
		provider := query.Get("provider")

		namespace := params["namespace"]

		modules, _ := r.ListModules(namespace, "", provider, offset, limit)

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
