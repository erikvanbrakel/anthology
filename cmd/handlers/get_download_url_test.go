package handlers_test

import (
	"testing"
	"github.com/gavv/httpexpect"
	"net/http"
)

func TestGetDownloadUrl(t *testing.T) {
	e := httpexpect.New(t, server.URL)

	r := e.GET("/v1/modules/namespace1/module1/provider1/1.0.0/download").Expect().Status(http.StatusNoContent)
	r.Body().Empty()
	r.Header("X-Terraform-Get").NotEmpty()

	e.GET("/v1/modules/namespace1/module1/provider1/5.0.0/download").Expect().Status(http.StatusNotFound)
}