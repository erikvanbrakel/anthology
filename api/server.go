package api

import (
	"net/http"
	"log"
	"github.com/spf13/viper"
	"fmt"
	"os"
	"os/signal"
	"context"
	"github.com/erikvanbrakel/anthology/api/v1"
	"github.com/erikvanbrakel/anthology/services"
	"github.com/erikvanbrakel/anthology/registry"
	"github.com/go-chi/chi"
	"errors"
	"strings"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	*http.Server
}

func ConfigureRegistry() (registry.Registry, error) {

	switch backend := viper.Get("backend"); backend {

	case "s3":
		log.Println("Configuring S3 backend")
		bucket := viper.GetString("s3.bucket")
		endpoint := viper.GetString("s3.endpoint")
		return registry.NewS3Registry(bucket, endpoint)

	case "filesystem":
		log.Println("Configuring Filesystem backend")
		basepath := viper.GetString("filesystem.basepath")

		return registry.NewFilesystemRegistry(basepath)

	case "demo":
		log.Println("Configuring demo backend (in-memory)")

		r := registry.NewFakeRegistry()

		providers := []string{"aws", "azure", "gcp"}
		modules := []string{"platform", "application", "logs"}
		versions := []string{"1.0.3", "2.0.0", "2.0.1", "5.0.0"}
		for _, p := range providers {
			for _, m := range modules {
				for _, v := range versions {
					r.PublishModule("demo", m, p, v, strings.NewReader("dummy"))
				}
			}
		}

		return r, nil
	default:
		log.Printf("Unknown backend: %v", backend)
		return nil, errors.New("unknown backend")
	}
}

func NewServer() (*Server, error) {
	log.Println("Configuring the server")

	port := viper.GetInt("port")

	r, err := ConfigureRegistry()
	if err != nil {
		return nil, err
	}

	service, err  := services.NewModuleService(r)
	if err != nil {
		return nil, err
	}

	api, err := v1.NewAPI(service)
	if err != nil {
		return nil, err
	}

	rootHandler := chi.NewRouter()
	rootHandler.Use(middleware.Logger)
	rootHandler.Get("/.well-known/terraform.json", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{
			"modules.v1": "/v1",
		})
	})

	rootHandler.Mount("/v1", api.Router())


	return &Server{
		&http.Server{
			Addr: fmt.Sprintf(":%v", port),
			Handler: rootHandler,
		},
	}, nil
}

func (srv *Server) Start() {
	log.Println("starting server...")
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Server gracefully stopped")

}