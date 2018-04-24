package handlers_test

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestGetDownloadUrl(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(st *testing.T) {

				e := httpexpect.New(st, s.Server.URL)

				namespace, _ := uuid.NewV4()

				// Set up (publish the module)
				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/1.2.3", namespace.String())).
					Expect().Status(http.StatusOK)

				// assert success path
				r := e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1/1.2.3/download", namespace.String())).
					Expect().Status(http.StatusNoContent)
				r.Body().Empty()
				r.Header("X-Terraform-Get").NotEmpty()

				// assert failure path
				e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1/5.0.0/download", namespace.String())).
					Expect().Status(http.StatusNotFound)
			},
		)
	}
}
