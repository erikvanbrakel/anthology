package v1_test

import (
	"net/http/httptest"
	"testing"

	"bytes"
	"github.com/erikvanbrakel/anthology/api/v1"
	"github.com/erikvanbrakel/anthology/registry"
	"github.com/erikvanbrakel/anthology/services"
	"github.com/gavv/httpexpect"
	"net/http"
	"github.com/go-chi/chi"
	"math/rand"
	"time"
)

type apiTestCase struct {
	tag    string
	method string
	url    string
	body   string
	status int
	assert func(*testing.T, *httpexpect.Response, *httptest.Server)
}

func randomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func runAPITests(t *testing.T, dataset []testModule, tests []apiTestCase) {
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			r := registry.NewFakeRegistry()

			for _, m := range dataset {
				r.PublishModule(m.namespace, m.name, m.provider, m.version, bytes.NewBuffer(m.data))
			}

			s, _ := services.NewModuleService(r)
			root := chi.NewRouter()
			api,_ := v1.NewAPI(s)

			// use a random prefix to make sure relative URLs don't depend on absolute paths
			mountpath := "/" + randomString(3)

			root.Mount(mountpath, api.Router())

			server := httptest.NewServer(root)
			defer server.Close()

			e := httpexpect.WithConfig(httpexpect.Config{
				BaseURL: server.URL,
				Printers: []httpexpect.Printer {
					httpexpect.NewDebugPrinter(t, true),
				},
				Reporter: httpexpect.NewAssertReporter(t),
			})

			result := e.Request(test.method, mountpath + test.url).
				WithClient(&http.Client{
					CheckRedirect: func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					},
				}).
				WithHeader("Content-Type", "application/json").
				WithBytes([]byte(test.body)).
				Expect().Status(test.status)

			test.assert(t, result, server)
		})
	}
}

type testModule struct {
	namespace string
	name      string
	provider  string
	version   string
	data      []byte
}
