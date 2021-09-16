package v1_test

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProvider(t *testing.T) {
	dataset := []testProvider{
		{"namespace1", "provider1", "1.0.0", "linux", "amd64", nil},
		{"namespace1", "provider1", "1.0.0", "linux", "386", nil},
		{"namespace1", "provider1", "1.0.0", "darwin", "amd64", nil},
		{"namespace2", "provider2", "1.0.0", "linux", "amd64", nil},
		{"namespace2", "provider2", "2.0.0", "linux", "amd64", nil},
	}

	runProviderAPITests(t, dataset, []apiTestCase{
		{
			"get all providers",
			"GET", "/", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.Value("meta").Object().NotEmpty()
				result.Value("providers").Array().NotEmpty()
			},
		},

		{
			"get all providers for namespace",
			"GET", "/namespace1", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.Value("meta").Object().NotEmpty()

				providers := result.Value("providers").Array()
				providers.Length().Equal(3)

				for _, m := range providers.Iter() {
					m.Object().ValueEqual("namespace", "namespace1")
				}
			},
		},

		{
			"get all providers for namespace (not-exist)",
			"GET", "/absent-namespace", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestListProviderVersions(t *testing.T) {
	dataset := []testProvider{
		{"namespace1", "provider1", "1.0.0", "linux", "amd64", nil},
		{"namespace1", "provider1", "1.0.0", "linux", "386", nil},
		{"namespace1", "provider1", "1.0.0", "darwin", "amd64", nil},
		{"namespace2", "provider2", "1.0.0", "linux", "amd64", nil},
		{"namespace2", "provider2", "2.0.0", "linux", "amd64", nil},
	}

	runProviderAPITests(t, dataset, []apiTestCase{
		{
			"list available versions for a specific provider",
			"GET", "/namespace1/provider1/versions", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				versions := result.Value("versions").Array()

				versions.NotEmpty()

				for _, m := range versions.Iter() {
					m.Object().ValueEqual("namespace", "namespace1")
					m.Object().ValueEqual("name", "provider1")
					m.Object().Value("platforms").Array().Length().Equal(3)
				}
			},
		},
		{
			"list available versions for a specific provider (not-exist)",
			"GET", "/namespace1/absent-provider/versions", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestGetProviderDownloadUrl(t *testing.T) {
	dataset := []testProvider{
		{"namespace1", "provider1", "1.0.0", "linux", "amd64", nil},
		{"namespace1", "provider1", "1.0.0", "linux", "386", nil},
		{"namespace1", "provider1", "1.0.0", "darwin", "amd64", nil},
	}

	runProviderAPITests(t, dataset, []apiTestCase{
		{
			"download source code for a specific provider version",
			"GET", "/namespace1/provider1/1.0.0/download/linux/amd64", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {

				result := r.JSON().Object()

				result.ValueEqual("os", "linux")
				result.ValueEqual("arch", "amd64")
				result.ValueEqual("filename", "terraform-provider-provider1_1.0.0_linux_amd64.zip")
			},
		},
		{
			"download source code for a specific provider version (not-exist)",
			"GET", "/namespace1/provider1/4.0.0/download/linux/amd64", "",
			http.StatusNotFound,
			assertError("not found"),
		},
	})
}

func TestGetProvider(t *testing.T) {
	dataset := []testProvider{
		{"namespace1", "provider1", "1.0.0", "linux", "amd64", nil},
		{"namespace1", "provider1", "1.0.0", "linux", "386", nil},
		{"namespace1", "provider1", "1.0.0", "darwin", "amd64", nil},
		{"namespace2", "provider2", "1.0.0", "linux", "amd64", nil},
		{"namespace2", "provider2", "2.0.0", "linux", "amd64", nil},
	}

	runProviderAPITests(t, dataset, []apiTestCase{
		{
			"get a specific provider",
			"GET", "/namespace1/provider1/1.0.0/linux/amd64", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				r.JSON().Object().NotEmpty()
			},
		},
		{
			"get a specific provider (not-exist)",
			"GET", "/namespace1/providers1/13.0.0/linux/amd64", "",
			http.StatusNotFound,
			assertError("not found"),
		},
	})
}

func TestPublishProvider(t *testing.T) {
	dataset := []testProvider{}

	providerData := "some data"

	runProviderAPITests(t, dataset, []apiTestCase{
		{
			"publish a new provider",
			"POST", "/namespace1/provider1/3.0.0/linux/amd64", providerData,
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.ValueEqual("os", "linux")
				result.ValueEqual("arch", "amd64")
				result.ValueEqual("filename", "terraform-provider-provider1_3.0.0_linux_amd64.zip")
			},
		},
	})
}
