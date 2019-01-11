package v1

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/erikvanbrakel/anthology/models"
	"io"
	"net/http"
	"strconv"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/erikvanbrakel/anthology/services"
	"errors"
	"context"
	"github.com/spf13/viper"
)

type (
	moduleResource struct {
		services.ModuleService
	}
)

type API struct {
	Modules *moduleResource
}

func NewAPI(service services.ModuleService) (*API, error) {

	if service == nil {
		return nil, errors.New("service can't be nil")
	}

	return &API{
		Modules: &moduleResource{service },
	}, nil
}

func (a *API) Router() *chi.Mux {

	router := chi.NewRouter()
	rg := router.With(a.MetaCtx)

	r := a.Modules

		// List modules
	rg.
		With(a.ListCtx).
		Get("/", r.query)

	rg.
		With(a.ListCtx).
		Get("/{namespace}", r.query)

	// Search modules
	// rg.Get("/search")

	// List available versions for a specific module
	rg.Get("/{namespace}/{name}/{provider}/versions", r.queryVersions)

	// Download source code for a specific module version
	rg.With(a.ModuleCtx).Get("/{namespace}/{name}/{provider}/{version}/download", r.getDownloadUrl)

	// Download the latest version of a module
	rg.Get("/{namespace}/{name}/{provider}/download", r.getLatestDownloadUrl)

	// List latest version of module for all providers
	rg.Get("/{namespace}/{name}", r.queryLatest)

	// Latest version for a specific module provider
	rg.Get("/{namespace}/{name}/{provider}", r.getLatest)

	// Get a specific module
	rg.With(a.ModuleCtx).Get("/{namespace}/{name}/{provider}/{version}", r.get)

	if viper.GetBool("publishing.enabled") {
		// Publish a specific module
		rg.Post("/{namespace}/{name}/{provider}/{version}", r.publish)
	}

	rg.With(a.ModuleCtx).Get("/{namespace}/{name}/{provider}/{version}/data.tgz", r.getModuleData)

	return router
}

func (a *API) ListCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		page := &PaginationInfo{
			CurrentOffset: 0,
			Limit: 10,
		}

		if o := r.URL.Query().Get("offset"); o != "" {
			offset, err := strconv.Atoi(o)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				return
			}
			page.CurrentOffset = offset
		}

		if lim := r.URL.Query().Get("offset"); lim != "" {
			limit, err := strconv.Atoi(lim)
			if err != nil {
				render.Status(r, http.StatusBadRequest)
				return
			}
			page.Limit = limit
		}

		ctx := context.WithValue(r.Context(), "pagination", page)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) MetaCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		metadata := &ModuleMetadata{
			Namespace: chi.URLParam(r, "namespace"),
			Name: chi.URLParam(r, "name"),
			Provider: chi.URLParam(r, "provider"),
			Version: chi.URLParam(r, "version"),
		}

		ctx := context.WithValue(r.Context(), "module_metadata", metadata)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) ModuleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		m := r.Context().Value("module_metadata").(*ModuleMetadata)
		module, err := a.Modules.Get(m.Namespace, m.Name, m.Provider, m.Version)

		if err != nil {
			render.Render(w,r,ErrInternalServerError(err))
			return
		}

		if module == nil {
			render.Render(w,r,ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), "module", module)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type ModuleMetadata struct {
	Namespace string
	Name string
	Provider string
	Version string
}

func (m *moduleResource) getModuleData(w http.ResponseWriter, r *http.Request) {

	module := r.Context().Value("module").(*models.Module)

	data, err := module.Data()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		_, err = io.Copy(w, data)
	}

}

func (m *moduleResource) publish(w http.ResponseWriter, r *http.Request) {

	meta := r.Context().Value("module_metadata").(*ModuleMetadata)

	maximumSize := viper.GetInt64("publishing.maximum_size") * 1024
	if r.ContentLength > maximumSize {
		render.Render(w, r, ErrPayloadTooLarge())
		return
	}

	err := m.Publish(meta.Namespace, meta.Name, meta.Provider, meta.Version, io.LimitReader(r.Body, r.ContentLength))
	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	m.getDownloadUrl(w, r)
}

func (m *moduleResource) query(w http.ResponseWriter, r *http.Request) {

	namespace :=  chi.URLParam(r, "namespace")
	provider := chi.URLParam(r, "provider")
	verified, _ := strconv.ParseBool(chi.URLParam(r, "verified"))

	page := r.Context().Value("pagination").(*PaginationInfo)

	modules, count, err := m.Query(namespace, "", provider, verified, page.CurrentOffset, page.Limit)

	if err != nil {
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	if count == 0 {
		render.Render(w,r,ErrNotFound())
		return
	}

	render.JSON(w, r, PaginatedList{
		PaginationInfo: *page,
		Modules:        modules,
	})
}

func (m *moduleResource) queryVersions(w http.ResponseWriter, req *http.Request) {

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")
	provider := chi.URLParam(req, "provider")

	versionsByModule, err := m.QueryVersions(namespace, name, provider)

	if err != nil {
		render.Render(w, req, ErrInternalServerError(err))
		return
	}

	if len(versionsByModule) == 0 {
		render.Render(w, req, ErrNotFound())
		return
	}

	render.JSON(w, req, struct {
		Modules VersionsList `json:"modules"`
	}{
		VersionsList{
			{
				Source:   fmt.Sprintf("%s/%s/%s", namespace, name, provider),
				Versions: versionsByModule,
			},
		},
	})
}

func (m *moduleResource) getDownloadUrl(w http.ResponseWriter, r *http.Request) {

	meta := r.Context().Value("module_metadata").(*ModuleMetadata)

	if exists, _ := m.Exists(meta.Namespace, meta.Name, meta.Provider, meta.Version); exists {
		render.Status(r, http.StatusNoContent)
		render.Render(w, r, NewTerraformModuleUrl(meta))

		return
	}

	render.Render(w, r, ErrNotFound())
}

type TerraformModuleUrl struct {
	*ModuleMetadata
}

func NewTerraformModuleUrl(metadata *ModuleMetadata) render.Renderer {
	return &TerraformModuleUrl { metadata }
}

func (r *TerraformModuleUrl) Render(w http.ResponseWriter, req *http.Request) error {
	rc := chi.RouteContext(req.Context())
	prefix := req.URL.Path[: (len(req.URL.Path) - len(rc.RoutePath))]

	w.Header().Set("X-Terraform-Get", fmt.Sprintf("%v/%v/%v/%v/%v/data.tgz", prefix, r.Namespace, r.Name, r.Provider, r.Version ))
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (m *moduleResource) getLatestDownloadUrl(w http.ResponseWriter, req *http.Request) {

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")
	provider := chi.URLParam(req, "provider")

	modules, count, err := m.Query(namespace, name, provider, false, 0, 100000)

	if err != nil {
		render.Render(w, req, ErrInternalServerError(err))
		return
	}

	if count == 0 {
		render.Render(w, req, ErrNotFound())
		return
	}

	var latestVersion semver.Version

	for _, m := range modules {
		moduleVersion, _ := semver.Make(m.Version)
		if moduleVersion.Compare(latestVersion) > 0 {
			latestVersion = moduleVersion
		}
	}

	prefix := req.URL.Path[:len(req.URL.Path) - len("/download")]
	url := fmt.Sprintf("%v/%v/download", prefix, latestVersion)

	render.Render(w, req, &Redirect{ URL: url })
}

type Redirect struct {
	URL string
}

func (r *Redirect) Render(w http.ResponseWriter, req *http.Request) error {
	w.Header().Set("Location", r.URL)
	render.Status(req, http.StatusFound)
	return nil
}

func (m *moduleResource) get(w http.ResponseWriter, req *http.Request) {

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")
	provider := chi.URLParam(req, "provider")
	version := chi.URLParam(req, "version")

	module, err := m.Get(namespace, name, provider, version)

	if err != nil {
		render.Render(w, req, ErrInternalServerError(err))
		return
	}

	if module == nil {
		render.Render(w, req, ErrNotFound())
		return
	}

	render.JSON(w, req, module)
}

func (m *moduleResource) getLatest(w http.ResponseWriter, req *http.Request) {

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")
	provider := chi.URLParam(req, "provider")

	modules, err := m.QueryVersions(namespace, name, provider)

	if err != nil {
		render.Render(w, req, ErrInternalServerError(err))
		return
	}

	if len(modules) == 0 {
		render.Render(w, req, ErrNotFound())
		return
	}

	var module = models.Module{
		Version: "0.0.0",
	}

	for _, m := range modules {
		moduleVersion, _ := semver.Make(m.Version)
		latestVersion, _ := semver.Make(module.Version)

		if moduleVersion.Compare(latestVersion) > 0 {
			module = m
		}
	}

	render.JSON(w, req, module)
}

func (m *moduleResource) queryLatest(w http.ResponseWriter, req *http.Request) {

	namespace := chi.URLParam(req, "namespace")
	name := chi.URLParam(req, "name")

	modules, err := m.QueryVersions(namespace, name, "")

	if err != nil {
		render.Render(w, req, ErrInternalServerError(err))
		return
	}

	if len(modules) == 0 {
		render.Render(w, req, ErrNotFound())
		return
	}

	var latestVersions = map[string]models.Module{}

	for _, m := range modules {
		latestVersion := latestVersions[m.Provider]
		latest, _ := semver.Make(latestVersion.Version)
		current, _ := semver.Make(m.Version)

		if current.Compare(latest) > 0 {
			latestVersions[m.Provider] = m
		}
	}

	v := make([]models.Module, 0, len(latestVersions))

	for _, value := range latestVersions {
		v = append(v, value)
	}

	render.JSON(w, req, PaginatedList{
		PaginationInfo: getPaginationInfo(chi.RouteContext(req.Context()), len(v)),
		Modules:        v,
	})
}

type VersionsList []struct {
	Source   string          `json:"source"`
	Versions []models.Module `json:"versions"`
}

func getPaginationInfo(c *chi.Context, count int) PaginationInfo {

	limit,_ := strconv.Atoi(c.URLParam("limit"))
	offset,_ := strconv.Atoi(c.URLParam("offset"))

	return PaginationInfo{
		CurrentOffset: offset,
		Limit:         limit,
	}
}

type PaginatedList struct {
	PaginationInfo PaginationInfo `json:"meta"`
	Modules        interface{}    `json:"modules"`
}

type PaginationInfo struct {
	Limit          int    `json:"limit"`
	PreviousOffset int    `json:"previous_offset"`
	PreviousUrl    int    `json:"previous_url"`
	CurrentOffset  int    `json:"current_offset"`
	NextOffset     int    `json:"next_offset"`
	NextUrl        string `json:"next_url"`
}
