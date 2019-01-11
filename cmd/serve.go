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
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringP("port", "", "8080", "Port the service listens on")
	serveCmd.Flags().StringP("backend", "", "demo", "Which backend to use. Valid options are s3, filesystem or demo.")

	serveCmd.Flags().StringP("ssl.certificate", "", "", "")
	serveCmd.Flags().StringP("ssl.key", "", "", "")

	serveCmd.Flags().StringP("filesystem.basepath", "", "", "")

	serveCmd.Flags().StringP("s3.bucket", "", "", "")
	serveCmd.Flags().StringP("s3.endpoint", "", "", "")
	serveCmd.Flags().StringP("s3.region", "", "us-east-1", "")

	serveCmd.Flags().BoolP("publishing.enabled","",false, "")
	serveCmd.Flags().StringP("publishing.maximum_size","","1000", "")
	viper.BindPFlags(serveCmd.Flags())
}
