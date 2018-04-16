package cmd

import (
	"github.com/erikvanbrakel/terraform-registry/cmd/registry"
	"github.com/erikvanbrakel/terraform-registry/cmd/handlers"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"fmt"
	"net/http"
	"path/filepath"
	"os"
)

type RegistryServerConfig struct {
	CertFile string
	KeyFile  string
	Port     int

	BasePath string
}

type RegistryServer struct {
	Router   *mux.Router
	Registry registry.Registry

	Port int
	CertFile, KeyFile string
}

func LoggingMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var f = logrus.Fields{}

		for k,v := range vars {
			f[k] = v
		}
		logrus.WithFields(f).Infof("%s - - [%s]", r.RemoteAddr, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func normalizePath (input string) string {
	r := filepath.FromSlash(input)

	if r[len(r)-1] != os.PathSeparator {
		r = r + string(os.PathSeparator)
	}

	return r
}

func NewServer(config RegistryServerConfig) (*RegistryServer, error) {

	router := mux.NewRouter()
	var r registry.Registry

	r = &registry.FilesystemRegistry{BasePath: normalizePath(config.BasePath)}

	router.Use(LoggingMiddleware)
	router.HandleFunc("/.well-known/terraform.json", handlers.ServiceDiscoveryHandler()).Methods("GET")

	v1 := router.PathPrefix("/v1/").Subrouter()
	v1.HandleFunc("/download/{namespace}/{name}/{provider}/{version}.tgz", handlers.DownloadHandler(config.BasePath)).Methods("GET")

	api := v1.PathPrefix("/modules/").Subrouter()

	api.HandleFunc("/", handlers.ListModulesHandler(r)).Methods("GET")
	api.HandleFunc("/{namespace}", handlers.ListModulesHandler(r)).Methods("GET")

	api.HandleFunc("/search", handlers.SearchModulesHandler()).Methods("GET")

	api.HandleFunc("/{namespace}/{name}", handlers.ListVersionsHandler(r)).Methods("GET")
	api.HandleFunc("/{namespace}/{name}/{provider}/versions", handlers.ListVersionsHandler(r)).Methods("GET")

	api.HandleFunc("/{namespace}/{name}/{provider}/download", handlers.GetDownloadUrlHandler(r)).Methods("GET")
	api.HandleFunc("/{namespace}/{name}/{provider}/{version}/download", handlers.GetDownloadUrlHandler(r)).Methods("GET")

	api.HandleFunc("/{namespace}/{name}/{provider}/{version}", handlers.GetModuleHandler(r)).Methods("GET")

	api.HandleFunc("/{namespace}/{name}/{provider}", handlers.GetModuleHandler(r)).Methods("GET")

	return &RegistryServer{
		Router: router,
		Port: config.Port,
		CertFile: config.CertFile,
		KeyFile: config.KeyFile,
	}, nil
}

func (server *RegistryServer) Run() {
	if server.CertFile != "" && server.KeyFile != "" {
		logrus.Infof("Starting registry. Listening on https://*:%d", server.Port)
		logrus.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", server.Port), server.CertFile, server.KeyFile, server.Router))
	} else {
		logrus.Infof("Starting registry. Listening on http://*:%d", server.Port)
		logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.Port), server.Router))
	}
}
