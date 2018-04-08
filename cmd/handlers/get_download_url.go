package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
)

func GetDownloadUrlHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		params := mux.Vars(request)

		namespace, name, provider, version := params["namespace"], params["name"], params["provider"], params["version"]

		if r.ModuleExists(namespace, name, provider, version) {
			writer.Header().Add("X-Terraform-Get", fmt.Sprintf("/v1/download/%s/%s/%s/%s.tgz", namespace, name, provider, version))
			writer.WriteHeader(http.StatusNoContent)
		} else {
			writer.WriteHeader(http.StatusNotFound)
		}
	}
}
