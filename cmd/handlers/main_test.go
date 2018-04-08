package handlers_test

import (
	"net/http/httptest"
	"testing"
	"github.com/erikvanbrakel/terraform-registry/cmd"
	"os"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	s,_ := cmd.NewServer(cmd.RegistryServerConfig {
		BasePath:"../../test/modules",
	})


	server = httptest.NewServer(s.Router)

	code := m.Run()

	os.Exit(code)
}
