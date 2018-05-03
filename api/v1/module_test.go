package v1_test

import (
	"github.com/gavv/httpexpect"
	"net/http"
	"net/http/httptest"
	"testing"
)

var errorNotFound = "not found"

func TestListModule(t *testing.T) {
	dataset := []testModule{
		{"namespace1", "module1", "aws", "1.0.0", nil},
		{"namespace1", "module1", "azure", "1.0.0", nil},
		{"namespace1", "module1", "gcp", "1.0.0", nil},
		{"namespace2", "module1", "aws", "1.0.0", nil},
		{"namespace2", "module1", "aws", "2.0.0", nil},
	}

	runAPITests(t, dataset, []apiTestCase{
		{
			"get all modules",
			"GET", "/", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.Value("meta").Object().NotEmpty()
				result.Value("modules").Array().NotEmpty()
			},
		},

		{
			"get all modules for namespace",
			"GET", "/namespace1", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.Value("meta").Object().NotEmpty()

				modules := result.Value("modules").Array()
				modules.Length().Equal(3)

				for _, m := range modules.Iter() {
					m.Object().ValueEqual("namespace", "namespace1")
				}
			},
		},

		{
			"get all modules for namespace (not-exist)",
			"GET", "/absent-namespace", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestListModuleVersions(t *testing.T) {
	dataset := []testModule{
		{"namespace1", "module1", "aws", "1.0.0", nil},
		{"namespace1", "module1", "aws", "2.0.0", nil},
		{"namespace2", "module1", "azure", "1.0.0", nil},
		{"namespace2", "module1", "gcp", "1.0.0", nil},
	}

	runAPITests(t, dataset, []apiTestCase{
		{
			"list available versions for a specific module",
			"GET", "/namespace1/module1/aws/versions", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				modules := result.Value("modules").Array()

				modules.NotEmpty()

				for _, m := range modules.Iter() {
					m.Object().ValueEqual("source", "namespace1/module1/aws")
					m.Object().Value("versions").Array().Length().Equal(2)
				}
			},
		},
		{
			"list available versions for a specific module (not-exist)",
			"GET", "/namespace1/absent-module/aws/versions", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestGetDownloadUrl(t *testing.T) {
	dataset := []testModule{
		{"namespace1", "module1", "aws", "1.0.0", nil},
		{"namespace1", "module1", "aws", "2.0.0", nil},
		{"namespace1", "module1", "aws", "3.0.0", nil},
	}

	runAPITests(t, dataset, []apiTestCase{
		{
			"download source code for a specific module version",
			"GET", "/namespace1/module1/aws/1.0.0/download", "",
			http.StatusNoContent,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				r.Header("X-Terraform-Get").NotEmpty()
			},
		},
		{
			"download source code for a specific module version (not-exist)",
			"GET", "/namespace1/module1/aws/4.0.0/download", "",
			http.StatusNotFound,
			assertError("not found"),
		},
		{
			"download the latest version of a module",
			"GET", "/namespace1/module1/aws/download", "",
			http.StatusFound,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				r.Header("Location").NotEmpty()
			},
		},
		{
			"download the latest version of a module (not-exist)",
			"GET", "/namespace1/module1/gcp/download", "",
			http.StatusNotFound,
			assertError("not found"),
		},
	})
}

func TestListLatestVersions(t *testing.T) {
	dataset := []testModule{
		{"namespace1", "module1", "aws", "1.0.0", nil},
		{"namespace1", "module1", "aws", "2.0.0", nil},
		{"namespace1", "module1", "gcp", "4.0.0", nil},
		{"namespace1", "module1", "gcp", "3.0.0", nil},
		{"namespace1", "module1", "azure", "7.0.0", nil},
		{"namespace1", "module1", "azure", "6.6.0", nil},
	}

	runAPITests(t, dataset, []apiTestCase{
		{
			"list latest version of module for all providers",
			"GET", "/namespace1/module1", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				r.JSON().Object().Keys().ContainsOnly("meta", "modules")

				for _, m := range r.JSON().Path("$.modules").Array().NotEmpty().Iter() {
					m.Object().ValueEqual("namespace", "namespace1")
					m.Object().ValueEqual("name", "module1")
				}
			},
		},
		{
			"list latest version of module for all providers (not-exist)",
			"GET", "/namespace1/absent-module", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestGetModule(t *testing.T) {
	dataset := []testModule{
		{"namespace1", "module1", "aws", "1.0.0", nil},
		{"namespace1", "module1", "aws", "2.0.0", nil},
		{"namespace1", "module1", "gcp", "4.0.0", nil},
		{"namespace1", "module1", "gcp", "3.0.0", nil},
		{"namespace1", "module1", "azure", "7.0.0", nil},
		{"namespace1", "module1", "azure", "6.6.0", nil},
	}

	runAPITests(t, dataset, []apiTestCase{
		{
			"get a specific module",
			"GET", "/namespace1/module1/gcp/3.0.0", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				r.JSON().Object().NotEmpty()
			},
		},
		{
			"get a specific module (not-exist)",
			"GET", "/namespace1/modules1/gcp/13.0.0", "",
			http.StatusNotFound,
			assertError("not found"),
		},
		{
			"latest version for a specific module provider",
			"GET", "/namespace1/module1/azure", "",
			http.StatusOK,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				result := r.JSON().Object()

				result.NotEmpty()
				result.ValueEqual("namespace", "namespace1")
				result.ValueEqual("name", "module1")
				result.ValueEqual("provider", "azure")
				result.ValueEqual("version", "7.0.0")
			},
		},

		{
			"latest version for a specific module provider (not-exist)",
			"GET", "/namespace1/absent-module/azure", "",
			http.StatusNotFound,
			assertError(errorNotFound),
		},
	})
}

func TestPublishModule(t *testing.T) {
	dataset := []testModule{}

	moduleData := "some data"

	runAPITests(t, dataset, []apiTestCase{
		{
			"publish a new module",
			"POST", "/namespace1/module1/gcp/3.0.0", moduleData,
			http.StatusNoContent,
			func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
				e := httpexpect.New(t, server.URL)

				downloadUrl := r.Header("X-Terraform-Get").NotEmpty()
				response := e.GET(downloadUrl.Raw()).Expect().Status(http.StatusOK)

				response.Body().Equal(moduleData)
			},
		},
	})
}

func assertError(error string) func(*testing.T, *httpexpect.Response, *httptest.Server) {
	return func(t *testing.T, r *httpexpect.Response, server *httptest.Server) {
		errors := r.JSON().Object().Value("errors").Array()
		errors.Contains(error)
	}
}
