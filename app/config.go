package app

import (
	"github.com/jessevdk/go-flags"
	"os"
)

var Config = &CommonOptions{}

type CommonOptions struct {
	Port       int               `short:"p" long:"port" description:"Port the service listens on" default:"8080"`
	Backend    string            `short:"b" long:"backend" choice:"s3" choice:"filesystem"`
	S3         S3Options         `group:"S3 configuration" namespace:"s3"`
	FileSystem FileSystemOptions `group:"Filesystem configuration" namespace:"filesystem"`
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
