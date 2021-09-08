package v1

import (
	"fmt"
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/go-ozzo/ozzo-routing"
	"io"
	"net/http"
	"strconv"
)

type (
	providerService interface {
		Query(rs app.RequestScope, namespace, name string, offset, limit int) ([]models.Provider, int, error)
		QueryVersions(rs app.RequestScope, namespace, name string) ([]models.ProviderVersions, error)
		Exists(rs app.RequestScope, namespace, name, version, OS, arch string) (bool, error)
		Get(rs app.RequestScope, namespace, name, version, OS, arch string) (*models.Provider, error)
		GetMetaData(rs app.RequestScope, namespace, name, version, OS, arch string) (models.ProviderDownload, error)
		GetData(rs app.RequestScope, namespace, name, version, OS, arch, file string) (io.Reader, error)
		Publish(rs app.RequestScope, namespace, name, version, OS, arch string, data io.Reader) error
	}

	providerResource struct {
		service providerService
	}
)

func ServeProviderResource(rg *routing.RouteGroup, service providerService) {
	r := &providerResource{service}

	// List providers
	rg.Get("/", r.query)
	rg.Get("/<namespace>", r.query)

	// Search providers
	rg.Get("/search")

	// List available versions for a specific provider
	rg.Get("/<namespace>/<name>/versions", r.queryVersions)

	// Download document for a specific provider version
	rg.Get("/<namespace>/<name>/<version>/download/<OS>/<arch>", r.getDownloadMetaData)

	// Download the requested file (not Hashicorp Registry endpoint)
	rg.Get("/<namespace>/<name>/<version>/download/<OS>/<arch>/<file>", r.getProviderData).Name("GetProviderData")

	// Get a specific provider
	rg.Get("/<namespace>/<name>/<version>/<OS>/<arch>", r.get)

	// Publish a specific provider
	rg.Post("/<namespace>/<name>/<version>/<OS>/<arch>", r.publish)

}

func (r *providerResource) getProviderData(c *routing.Context) error {
	rs := app.GetRequestScope(c)

	namespace := c.Param("namespace")
	name := c.Param("name")
	version := c.Param("version")
	OS := c.Param("OS")
	arch := c.Param("arch")
	file := c.Param("file")

	data, err := r.service.GetData(rs, namespace, name, version, OS, arch, file)

	fmt.Println("Received Data from Artifactory")
	if err != nil {
		return err
	}

	c.Response.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(c.Response, data)

	return err
}

func (r *providerResource) publish(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	namespace, name, version, OS, arch := c.Param("namespace"), c.Param("name"), c.Param("version"), c.Param("OS"), c.Param("arch")

	err := r.service.Publish(rs, namespace, name, version, OS, arch, c.Request.Body)
	if err != nil {
		return err
	}

	return r.getDownloadMetaData(c)
}

func (r *providerResource) query(c *routing.Context) error {
	rs := app.GetRequestScope(c)

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	namespace := c.Param("namespace")

	providers, count, err := r.service.Query(rs, namespace, "", offset, limit)

	if err != nil {
		return err
	}

	if count == 0 {
		c.Response.WriteHeader(http.StatusNotFound)
		return c.Write(apiError{[]string{"not found"}})
	}

	paginationInfo := getProviderPaginationInfo(c, count)

	return c.Write(ProviderPaginatedList{
		PaginationInfo: paginationInfo,
		Providers:      providers,
	})
}

func (r *providerResource) queryVersions(c *routing.Context) error {
	rs := app.GetRequestScope(c)

	namespace := c.Param("namespace")
	name := c.Param("name")

	versions, err := r.service.QueryVersions(rs, namespace, name)

	if err != nil {
		return err
	}

	if len(versions) == 0 {
		c.Response.WriteHeader(http.StatusNotFound)
		return c.Write(apiError{[]string{"not found"}})
	}

	return c.Write(
		ProviderVersionsList{
			Source:   fmt.Sprintf("%s/%s", namespace, name),
			Versions: versions,
		},
	)
}

func (r *providerResource) getDownloadMetaData(c *routing.Context) error {
	rs := app.GetRequestScope(c)

	namespace, name, version, OS, arch := c.Param("namespace"), c.Param("name"), c.Param("version"), c.Param("OS"), c.Param("arch")

	data, err := r.service.GetMetaData(rs, namespace, name, version, OS, arch)

	if err == nil {

		// Set URL paths, if they exist
		if data.DownloadURL != "" {
			data.DownloadURL = c.URL("GetProviderData",
				"namespace", namespace,
				"name", name,
				"version", version,
				"OS", OS,
				"arch", arch,
				"file", data.DownloadURL,
			)
		}
		if data.SHASumsURL != "" {
			data.SHASumsURL = c.URL("GetProviderData",
				"namespace", namespace,
				"name", name,
				"version", version,
				"OS", OS,
				"arch", arch,
				"file", data.SHASumsURL,
			)
		}
		if data.SHASumsSigURL != "" {
			data.SHASumsSigURL = c.URL("GetProviderData",
				"namespace", namespace,
				"name", name,
				"version", version,
				"OS", OS,
				"arch", arch,
				"file", data.SHASumsSigURL,
			)
		}
		c.Write(data)
		return nil
	}

	c.Response.WriteHeader(http.StatusNotFound)
	return c.Write(apiError{[]string{"not found"}})

}

func (r *providerResource) get(c *routing.Context) error {
	rs := app.GetRequestScope(c)

	namespace := c.Param("namespace")
	name := c.Param("name")
	version := c.Param("version")
	OS := c.Param("OS")
	arch := c.Param("arch")

	provider, err := r.service.Get(rs, namespace, name, version, OS, arch)

	if err != nil {
		return err
	}

	if provider == nil {
		c.Response.WriteHeader(http.StatusNotFound)
		return c.Write(apiError{[]string{"not found"}})
	}

	return c.Write(provider)
}

type ProviderVersionsList struct {
	Source   string                    `json:"source"`
	Versions []models.ProviderVersions `json:"versions"`
}

func getProviderPaginationInfo(c *routing.Context, count int) PaginationInfo {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	return PaginationInfo{
		CurrentOffset: offset,
		Limit:         limit,
	}
}

type ProviderPaginatedList struct {
	PaginationInfo PaginationInfo `json:"meta"`
	Providers      interface{}    `json:"providers"`
}
