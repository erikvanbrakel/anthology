package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"encoding/json"
	"github.com/erikvanbrakel/anthology/cmd/api"
)

func DownloadHandler(r registry.Registry) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		params := mux.Vars(request)

		namespace, name, provider, version := params["namespace"], params["name"], params["provider"], params["version"]

		reader, err := r.GetModuleData(namespace, name, provider, version)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(api.NewError(err.Error()))
			return
		} else {
			if reader == nil {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
		}
		writer.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s-%s-%s-%s.tgz", namespace, name, provider, version))

		written, _ := io.Copy(writer, reader)

		writer.Header().Set("Content-Length", strconv.FormatInt(written, 10))
	}
}

