package handlers_test

import (
	"testing"
	//"github.com/gavv/httpexpect"
	//"net/http"
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	"net/http"
)

func TestListModulesWithoutFilter(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(t *testing.T) {

				e := httpexpect.WithConfig(httpexpect.Config{
					BaseURL:  s.Server.URL,
					Reporter: httpexpect.NewAssertReporter(t),
					Printers: []httpexpect.Printer{
						httpexpect.NewDebugPrinter(t, true),
					},
				})
				namespace, _ := uuid.NewV4()

				e.POST(fmt.Sprintf("/v1/modules/%s/mod1/provider1/1.0.0", namespace.String())).
					Expect().Status(http.StatusOK)

				e.POST(fmt.Sprintf("/v1/modules/%s/mod1/provider1/2.0.0", namespace.String())).
					Expect().Status(http.StatusOK)

				json := e.GET("/v1/modules/").Expect().Status(http.StatusOK).JSON().Object()

				json.Keys().ContainsOnly("meta", "modules")

				json.Value("modules").Array().NotEmpty()
			},
		)
	}
}

func TestListModulesWithNamespace(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(t *testing.T) {
				e := httpexpect.New(t, s.Server.URL)

				namespace1, _ := uuid.NewV4()
				namespace2, _ := uuid.NewV4()

				e.POST(fmt.Sprintf("/v1/modules/%s/mod1/provider1/1.0.0", namespace1.String())).
					Expect().Status(http.StatusOK)
				e.POST(fmt.Sprintf("/v1/modules/%s/mod1/provider1/1.0.0", namespace2.String())).
					Expect().Status(http.StatusOK)

				json := e.GET(fmt.Sprintf("/v1/modules/%s", namespace1.String())).
					Expect().Status(http.StatusOK).JSON().Object()

				json.Keys().ContainsOnly("meta", "modules")

				for _, m := range json.Path("$.modules[*].namespace").Array().NotEmpty().Iter() {
					m.String().Equal(namespace1.String())
				}
			},
		)
	}
}
