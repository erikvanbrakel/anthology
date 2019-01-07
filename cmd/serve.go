package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/erikvanbrakel/anthology/api"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts an http server for the registry",
	Run: func(cmd *cobra.Command, args []string) {
		server, err := api.NewServer()
		if err != nil {
			log.Fatal(err)
		}
		server.Start()

/*		app.LoadConfig()

		logger := logrus.New()

		var r registry.Registry

		switch app.Config.Backend {
		case "s3":
			r = registry.NewS3Registry(app.Config.S3)
			break
		case "filesystem":
			r = registry.NewFilesystemRegistry(app.Config.FileSystem)
			break
		}
		http.Handle("/", buildRouter(logger, r))

		address := fmt.Sprintf(":%v", app.Config.Port)
		logger.Infof("server %v is started at %v", app.Version, address)

		if app.Config.SSLConfig.IsValid() {
			panic(http.ListenAndServeTLS(address, app.Config.SSLConfig.Certificate, app.Config.SSLConfig.Key, nil))
		} else {
			panic(http.ListenAndServe(address, nil))
		}*/
	},
}
/*
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
			"modules.v1": "/v1/modules/",
		})
	})

	rg := router.Group("/v1/modules")

	v1.ServeModuleResource(rg, services.NewModuleService(reg))
	return router
}*/

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "", "8080", "Port the service listens on")
	serveCmd.Flags().StringP("backend", "", "", "")
	serveCmd.Flags().StringP("ssl.certificate", "", "", "")
	serveCmd.Flags().StringP("ssl.key", "", "", "")
	serveCmd.Flags().StringP("filesystem.basepath", "", "", "")
	serveCmd.Flags().StringP("s3.bucket", "", "", "")
	serveCmd.Flags().StringP("s3.endpoint", "", "", "")

	viper.BindPFlags(serveCmd.Flags())
}
