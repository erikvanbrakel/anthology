package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/erikvanbrakel/anthology/cmd/api"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"github.com/gorilla/mux"
	"net/http"
)

func GetDownloadUrlHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {

	return func(writer http.ResponseWriter, request *http.Request) {

		params := mux.Vars(request)

		namespace, name, provider, version := params["namespace"], params["name"], params["provider"], params["version"]

		module, err := r.GetModule(namespace, name, provider, version)
		output := json.NewEncoder(writer)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			output.Encode(api.NewError(err.Error()))
		} else {
			if module != nil {
				writer.Header().Add("X-Terraform-Get", fmt.Sprintf("/v1/download/%s/%s/%s/%s.tgz", namespace, name, provider, version))
				writer.WriteHeader(http.StatusNoContent)
			} else {
				writer.WriteHeader(http.StatusNotFound)
				output.Encode(api.NewError("Not Found"))
			}
		}
	}
}
