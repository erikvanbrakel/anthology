package handlers

import (
	"encoding/json"
	"github.com/blang/semver"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/erikvanbrakel/anthology/cmd/api"
)

func GetModuleHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace, name, provider := params["namespace"], params["name"], params["provider"]

		version, hasVersion := params["version"]

		moduleVersions, _ := r.ListVersions(namespace,name,provider)

		var module *registry.Module

		for i, f := range moduleVersions[0].Versions {
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
						module = &moduleVersions[0].Versions[i]
					}
				} else {
					module = &moduleVersions[0].Versions[i]
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
