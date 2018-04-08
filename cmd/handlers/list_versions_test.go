package handlers_test

import (
	"testing"
	"net/http"
	"github.com/gavv/httpexpect"
)

func TestListAvailableVersionsForASpecificModule(t *testing.T) {

	e := httpexpect.New(t, server.URL)

	json := e.GET("/v1/modules/namespace1/module1/provider1/versions").Expect().Status(http.StatusOK).JSON().Object()

	json.Keys().ContainsOnly("modules")

	for _, m := range json.Path("$.modules[0].versions[*].id").Array().NotEmpty().Iter() {
		m.String().Contains("namespace1/module1/provider1/")
	}
}
