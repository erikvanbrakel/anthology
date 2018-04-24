package handlers_test

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestGetModule(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(t *testing.T) {

				e := httpexpect.New(t, s.Server.URL)

				namespace, _ := uuid.NewV4()

				// Set up (publish the module)
				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/3.2.1", namespace.String())).Expect().Status(http.StatusOK)

				// assert success path
				json := e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1/3.2.1", namespace.String())).Expect().Status(http.StatusOK).JSON().Object()

				json.Value("id").String().Equal(fmt.Sprintf("%s/module1/provider1/3.2.1", namespace.String()))

				// assert failure path
				e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1/5.0.0", namespace.String())).Expect().Status(http.StatusNotFound)
			},
		)
	}
}

func TestGetLatestModule(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(t *testing.T) {
				e := httpexpect.New(t, s.Server.URL)

				namespace, _ := uuid.NewV4()
				// Set up (publish the module)
				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/3.2.1", namespace.String())).Expect().Status(http.StatusOK)
				e.POST(fmt.Sprintf("/v1/modules/%s/module2/provider1/5.3.10", namespace.String())).Expect().Status(http.StatusOK)
				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/4.7.2", namespace.String())).Expect().Status(http.StatusOK)

				json := e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1", namespace.String())).Expect().Status(http.StatusOK).JSON().Object()

				json.Value("id").String().Equal(fmt.Sprintf("%s/module1/provider1/4.7.2", namespace.String()))
			},
		)
	}
}
