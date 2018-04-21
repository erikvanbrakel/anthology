package handlers_test

import (
	"net/http/httptest"
	"testing"
	"github.com/erikvanbrakel/anthology/cmd"
	"os"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"strings"
)

var server *httptest.Server

func TestMain(m *testing.M) {

	r := createFakeRegistry()

	s,_ := cmd.NewServer(cmd.RegistryServerConfig {
		BasePath:"../../test/modules",
		Bucket: "modules",
	}, r)


	server = httptest.NewServer(s.Router)

	code := m.Run()

	os.Exit(code)
}

func createFakeRegistry() *registry.FakeRegistry {
	r := registry.FakeRegistry{}

	for _, n := range []string{"namespace1", "namespace2", "namespace3"} {
		for _, m := range []string{"module1", "module2", "module3"} {
			for _, p := range []string{"provider1", "provider2", "provider3"} {
				for _, v := range []string{"1.0.0", "2.0.0", "3.0.0"} {
					r.Modules = append(r.Modules, registry.Module{
						Namespace: n,
						Name: m,
						Provider: p,
						Version: v,
						ID: strings.Join([]string { n,m,p,v },"/"),
					})
				}
			}
		}
	}

	return &r
}