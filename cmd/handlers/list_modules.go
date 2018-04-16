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

		modules,total,_ := r.ListModules(namespace, "", provider, offset, limit)

		previousOffset := offset - limit
		if previousOffset < 0 {
			previousOffset = 0
		}
		nextOffset := offset + limit
		if nextOffset > total {
			nextOffset = 0
		}

		currentRoute := mux.CurrentRoute(request)

		nextUrl := ""
		if nextOffset > 0 {
			nextRoute,_ := currentRoute.URL(
				"namespace",namespace,
				"provider",provider)

			q := nextRoute.Query()

			q.Set("offset", strconv.Itoa(nextOffset))
			q.Set("limit", strconv.Itoa(limit))
			q.Set("provider", provider)

			nextRoute.RawQuery = q.Encode()
			nextUrl = nextRoute.String()
		}

		previousUrl := ""

		if offset != 0 {
			previousRoute, _ := currentRoute.
				URL(
				"namespace", namespace,
				"provider", provider)
			q := previousRoute.Query()

			q.Set("offset", strconv.Itoa(previousOffset))
			q.Set("limit", strconv.Itoa(limit))
			q.Set("provider", provider)

			previousRoute.RawQuery = q.Encode()
			previousUrl = previousRoute.String()
		}
		meta := api.Meta{
			CurrentOffset: offset,
			PreviousOffset: previousOffset,
			Limit:         limit,
			NextOffset:    nextOffset,
			NextUrl: nextUrl,
			PreviousUrl: previousUrl,
		}

		writer.Header().Add("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(api.ListModulesResponse{Meta: meta, Modules: modules})
	}
}
