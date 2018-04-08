package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
	"github.com/erikvanbrakel/terraform-registry/cmd/api"
	"github.com/sirupsen/logrus"
)

func ListVersionsHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace, name, provider := params["namespace"], params["name"], params["provider"]

		logrus.Infof("ListVersions(namespace=%s,name=%s,provider=%s",namespace,name,provider)
		modules, _ := r.ListVersions(namespace, name, provider)
		logrus.Infof("Found %d versions of the module.", len(modules))
		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(api.ListVersionsResponse{Modules: modules})
	}
}
