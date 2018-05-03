package app

import (
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"os"
)

var Config = &CommonOptions{}

type CommonOptions struct {
	Port       int               `short:"p" long:"port" description:"Port the service listens on" default:"8080"`
	Backend    string            `short:"b" long:"backend" choice:"s3" choice:"filesystem"`
	S3         S3Options         `group:"S3 configuration" namespace:"s3"`
	FileSystem FileSystemOptions `group:"Filesystem configuration" namespace:"filesystem"`
	SSLConfig  SSLOptions        `group:"SSL Configuration" namespace:"ssl"`
}

type SSLOptions struct {
	Certificate string `long:"certificate" description:"Path to the SSL certificate"`
	Key         string `long:"key" description:"Path to the SSL certificate key"`
}

type S3Options struct {
	Bucket   string `long:"bucket" description:"S3 bucket to use as backing storage"`
	Endpoint string `long:"endpoint" description:"S3 endpoint"`
}

type FileSystemOptions struct {
	BasePath string `long:"basepath" description:"Basepath to store modules"`
}

func LoadConfig() error {
	p := flags.NewParser(Config, flags.Default)

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	return nil
}

func (o SSLOptions) IsValid() bool {
	if o.Certificate == "" && o.Key == "" {
		return false
	}

	if _, err := os.Stat(o.Certificate); err != nil {
		logrus.Warnf("SSL configuration not valid, certificate file %s not found", o.Certificate)
		return false
	}

	if _, err := os.Stat(o.Key); err != nil {
		logrus.Warnf("SSL configuration not valid, key file %s not found", o.Key)
		return false
	}
	return true
}
