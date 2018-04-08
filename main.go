package main

import (
	"flag"
	"github.com/erikvanbrakel/terraform-registry/cmd"
)

func main() {

	config := cmd.RegistryServerConfig{}

	flag.StringVar(&config.CertFile, "tls_cert", "", "TLS certificate file")
	flag.StringVar(&config.KeyFile, "tls_key", "", "TLS certificate key")
	flag.IntVar(&config.Port, "port", 1234, "server port")
	flag.StringVar(&config.BasePath, "module_path", "", "Base path for module storage")

	flag.Parse()

	server, _ := cmd.NewServer(config)
	server.Run()
}
