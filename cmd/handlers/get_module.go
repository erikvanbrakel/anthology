package handlers

import (
	"encoding/json"
	"github.com/blang/semver"
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/terraform-registry/cmd/api"
)

func GetModuleHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace, name, provider := params["namespace"], params["name"], params["provider"]

		version, hasVersion := params["version"]

		modules, _ := r.ListModules(namespace, name, provider, 0, 9999)

		var module *registry.Module

		for _, f := range modules {
			if hasVersion {
				if f.Version == version {
					module = &f
					break
				}
			} else {
				if module != nil {
					fver, _ := semver.Make(f.Version)
					cver, _ := semver.Make(module.Version)

					if fver.Compare(cver) > 0 {
						module = &f
					}
				} else {
					module = &f
				}
			}
		}

		if module != nil {
			writer.Header().Add("Content-Type", "application/json")
			json.NewEncoder(writer).Encode(&module)
		} else {
			writer.WriteHeader(404)
			e := api.ApiError{
				Errors: []string{"Not Found"},
			}
			json.NewEncoder(writer).Encode(e)
		}
	}
}
