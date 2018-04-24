package handlers_test

import (
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/satori/go.uuid"
	"net/http"
	"testing"
)

func TestListAvailableVersionsForASpecificModule(t *testing.T) {
	for _, s := range servers {
		t.Run(
			s.Name, func(t *testing.T) {

				e := httpexpect.New(t, s.Server.URL)

				namespace, _ := uuid.NewV4()

				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/1.0.0", namespace)).Expect().Status(http.StatusOK)
				e.POST(fmt.Sprintf("/v1/modules/%s/module1/provider1/2.0.0", namespace)).Expect().Status(http.StatusOK)
				e.POST(fmt.Sprintf("/v1/modules/%s/module2/provider1/3.0.0", namespace)).Expect().Status(http.StatusOK)

				json := e.GET(fmt.Sprintf("/v1/modules/%s/module1/provider1/versions", namespace.String())).
					Expect().Status(http.StatusOK).JSON().Object()

				json.Keys().ContainsOnly("modules")

				for _, m := range json.Path("$.modules[0].versions[*].id").Array().NotEmpty().Iter() {
					m.String().Contains(fmt.Sprintf("%s/module1/provider1/", namespace.String()))
				}
			},
		)
	}
}
