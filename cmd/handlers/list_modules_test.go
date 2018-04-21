package handlers_test

import (
	"testing"
	//"github.com/gavv/httpexpect"
	//"net/http"
	"github.com/gavv/httpexpect"
	"net/http"
)


func TestListModulesWithoutFilter(t *testing.T) {

	e := httpexpect.New(t, server.URL)

	json := e.GET("/v1/modules/").Expect().Status(http.StatusOK).JSON().Object()

	json.Keys().ContainsOnly("meta", "modules")

	json.Value("modules").Array()
}

func TestListModulesWithNamespace(t *testing.T) {

	e := httpexpect.New(t, server.URL)

	json := e.GET("/v1/modules/namespace1").Expect().Status(http.StatusOK).JSON().Object()

	json.Keys().ContainsOnly("meta", "modules")

	for _, m := range json.Path("$.modules[*].namespace").Array().NotEmpty().Iter() {
		m.String().Equal("namespace1")
	}
}
