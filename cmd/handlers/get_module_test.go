package handlers_test

import (
	"testing"
	"github.com/gavv/httpexpect"
	"net/http"
)

func TestGetModule(t *testing.T) {
	e := httpexpect.New(t, server.URL)

	json := e.GET("/v1/modules/namespace1/module1/provider1/1.0.0").Expect().Status(http.StatusOK).JSON().Object()

	json.Value("id").String().Equal("namespace1/module1/provider1/1.0.0")

	e.GET("/v1/modules/namespace1/module1/provider1/5.0.0").Expect().Status(http.StatusNotFound)
}

func TestGetLatestModule(t *testing.T) {
	e := httpexpect.New(t, server.URL)

	json := e.GET("/v1/modules/namespace1/module1/provider1").Expect().Status(http.StatusOK).JSON().Object()

	json.Value("id").String().Equal("namespace1/module1/provider1/3.0.0")
}