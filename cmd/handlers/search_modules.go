package handlers

import "net/http"

func SearchModulesHandler() func(http.ResponseWriter, *http.Request) {
	return func (writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Content-Type", "application/json")
	}
}
