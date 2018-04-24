package main

import (
	"flag"
	"github.com/erikvanbrakel/anthology/cmd"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	config := cmd.RegistryServerConfig{}

	flag.StringVar(&config.CertFile, "tls_cert", "", "TLS certificate file")
	flag.StringVar(&config.KeyFile, "tls_key", "", "TLS certificate key")
	flag.IntVar(&config.Port, "port", 1234, "server port")

	var bucket string
	flag.StringVar(&bucket, "bucket", "", "Bucket name of s3 storage")
	flag.Parse()

	r := registry.S3Registry{Bucket: bucket}
	server, _ := cmd.NewServer(config, &r)
	go server.Run()

	var gracefulStop = make(chan os.Signal)

	// watch for SIGTERM and SIGINT from the operating system, and notify the app on
	// the gracefulStop channel
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	signal.Notify(gracefulStop, syscall.SIGKILL)

	<-gracefulStop

	os.Exit(0)
}
