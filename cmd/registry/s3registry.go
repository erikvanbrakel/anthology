package registry

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"strings"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/sirupsen/logrus"
)

type S3Registry struct {
	Bucket string
}

func (r *S3Registry) ListModules(namespace, name, provider string, offset, limit int) (modules []Module, total int, err error) {

	modules,err = r.getModules(namespace,name,provider)

	if err != nil {
		return nil, 0, err
	}

	return modules,len(modules),nil
}

func (r *S3Registry) ListVersions(namespace, name, provider string) (versions []ModuleVersions, err error) {
	modules,err := r.getModules(namespace, name, provider)

	versions = []ModuleVersions{
		{
			Source: strings.Join([]string { namespace, name, provider }, "/"),
			Versions: modules,
		},
	}

	return versions, err
}

func (r *S3Registry) GetModule(namespace, name, provider, version string) (module *Module, err error) {
	config := &aws.Config{
		DisableSSL: aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Region: aws.String("us-east-1"),
	}

	config.CredentialsChainVerboseErrors = aws.Bool(true)

	s,err := session.NewSession(config)

	if err != nil {
		return nil, err
	}

	s3client := s3.New(s)

	_,err = s3client.GetObject(&s3.GetObjectInput{
		Key: aws.String("/" + strings.Join([]string { namespace, name, provider, version }, "/") + ".tgz"),
		Bucket: aws.String(r.Bucket),
	})

	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return nil,nil
			}
			return nil, err
		}
	}

	return &Module {
		Namespace: namespace,
		Name: name,
		Provider: provider,
		Version: version,
	},nil
}

func (r *S3Registry) getModules(namespace,name,provider string) (modules []Module, err error) {
	config := &aws.Config{
		DisableSSL: aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Region: aws.String("us-east-1"),
	}
	config.CredentialsChainVerboseErrors = aws.Bool(true)

	prefix := ""

	if namespace != "" {
		prefix = namespace
		if name != "" {
			prefix = strings.Join([]string { prefix, name}, "/")
			if provider != "" {
				prefix = strings.Join([]string { prefix, provider}, "/")
			}
		}
	}

	s,err := session.NewSession(config)

	if err != nil {
		return nil, err
	}

	s3client := s3.New(s)

	loi := s3.ListObjectsInput{
		Bucket: aws.String(r.Bucket),
		Prefix: aws.String(prefix),
	}
	result,err := s3client.ListObjects(&loi)

	if err != nil {
		logrus.Errorf("error: %s", err.(awserr.Error))
		return nil, err
	}

	for _,o := range result.Contents {
		parts := strings.Split(*o.Key, "/")

		if len(parts) == 4 {
			modules = append(modules, Module {
				ID: strings.TrimRight(*o.Key, ".tgz"),
				Namespace: parts[0],
				Name: parts[1],
				Provider: parts[2],
				Version: strings.TrimRight(parts[3], ".tgz"),
			})
		}
	}

	return modules, nil
}