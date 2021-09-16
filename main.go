package main

import (
	"fmt"
	"github.com/erikvanbrakel/anthology/api/v1"
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/registry"
	"github.com/erikvanbrakel/anthology/services"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	if err := app.LoadConfig(); err != nil {
		panic(fmt.Errorf("invalid configuration: %s", err))
	}

	logger := logrus.New()

	var r registry.Registry

	switch app.Config.Backend {
	case "s3":
		r = registry.NewS3Registry(app.Config.S3)
		break
	case "filesystem":
		r = registry.NewFilesystemRegistry(app.Config.FileSystem)
		break
	case "artifactory":
		r = registry.NewArtifactoryRegistry(app.Config.Artifactory)
		break
	}
	http.Handle("/", buildRouter(logger, r))

	address := fmt.Sprintf(":%v", app.Config.Port)
	logger.Infof("server %v is started at %v", app.Version, address)

	if app.Config.SSLConfig.IsValid() {
		panic(http.ListenAndServeTLS(address, app.Config.SSLConfig.Certificate, app.Config.SSLConfig.Key, nil))
	} else {
		panic(http.ListenAndServe(address, nil))
	}
}

func buildRouter(logger *logrus.Logger, reg registry.Registry) *routing.Router {
	router := routing.New()

	router.To("GET,HEAD", "/ping", func(c *routing.Context) error {
		c.Abort()
		return c.Write("OK" + app.Version)
	})

	router.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),
	)

	router.To("GET", "/.well-known/terraform.json", func(c *routing.Context) error {
		c.Abort()
		return c.Write(map[string]string{
			"modules.v1":   "/v1/modules/",
			"providers.v1": "/v1/providers/",
		})
	})

	mrg := router.Group("/v1/modules")

	v1.ServeModuleResource(mrg, services.NewModuleService(reg))

	prg := router.Group("/v1/providers")

	v1.ServeProviderResource(prg, services.NewProviderService(reg))

	return router
}
