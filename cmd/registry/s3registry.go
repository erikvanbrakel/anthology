package registry

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type S3Registry struct {
	Bucket   string
	Endpoint string
}

func (r *S3Registry) ListModules(namespace, name, provider string, offset, limit int) (modules []Module, total int, err error) {

	modules, err = r.getModules(namespace, name, provider)

	if err != nil {
		return nil, 0, err
	}

	return modules, len(modules), nil
}

func (r *S3Registry) ListVersions(namespace, name, provider string) (versions []ModuleVersions, err error) {
	modules, err := r.getModules(namespace, name, provider)

	versions = []ModuleVersions{
		{
			Source:   strings.Join([]string{namespace, name, provider}, "/"),
			Versions: modules,
		},
	}

	return versions, err
}

func (r *S3Registry) GetModule(namespace, name, provider, version string) (module *Module, err error) {

	_, err = r.getObject("/" + strings.Join([]string{namespace, name, provider, version}, "/") + ".tgz")

	if err != nil {

		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return nil, nil
			}
			return nil, err
		}
	}

	return &Module{
		Namespace: namespace,
		Name:      name,
		Provider:  provider,
		Version:   version,
	}, nil
}

func (r *S3Registry) getObject(key string) (*s3.GetObjectOutput, error) {

	s3client := s3.New(r.getSession())

	return s3client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(r.Bucket),
	})

}

func (r *S3Registry) GetModuleData(namespace, name, provider, version string) (*bytes.Buffer, error) {
	object, err := r.getObject("/" + strings.Join([]string{namespace, name, provider, version}, "/") + ".tgz")

	buffer := &bytes.Buffer{}

	if err != nil {
		return nil, nil
	}

	io.Copy(buffer, object.Body)

	return buffer, nil
}

func (r *S3Registry) getSession() *session.Session {
	config := &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
	}
	if r.Endpoint != "" {
		if !strings.HasPrefix(r.Endpoint, "https") {
			config.DisableSSL = aws.Bool(true)
		}
		config.Endpoint = aws.String(r.Endpoint)
	}
	config.CredentialsChainVerboseErrors = aws.Bool(true)

	s, _ := session.NewSession(config)

	return s
}

func (r *S3Registry) PublishModule(namespace, name, provider, version string, data io.Reader) (err error) {

	key := strings.Join([]string{namespace, name, provider, version}, "/") + ".tgz"

	uploader := s3manager.NewUploader(r.getSession())

	input := s3manager.UploadInput{
		Body:   data,
		Key:    aws.String(key),
		Bucket: aws.String(r.Bucket),
	}

	_, err = uploader.Upload(&input)

	return err
}

func (r *S3Registry) getModules(namespace, name, provider string) (modules []Module, err error) {

	prefix := ""

	if namespace != "" {
		prefix = namespace
		if name != "" {
			prefix = strings.Join([]string{prefix, name}, "/")
			if provider != "" {
				prefix = strings.Join([]string{prefix, provider}, "/")
			}
		}
	}

	s3client := s3.New(r.getSession())

	loi := s3.ListObjectsInput{
		Bucket: aws.String(r.Bucket),
		Prefix: aws.String(prefix),
	}
	result, err := s3client.ListObjects(&loi)

	if err != nil {
		logrus.Errorf("error: %s", err.(awserr.Error))
		return nil, err
	}

	for _, o := range result.Contents {
		parts := strings.Split(*o.Key, "/")

		if len(parts) == 4 {
			modules = append(modules, Module{
				ID:        strings.TrimRight(*o.Key, ".tgz"),
				Namespace: parts[0],
				Name:      parts[1],
				Provider:  parts[2],
				Version:   strings.TrimRight(parts[3], ".tgz"),
			})
		}
	}

	return modules, nil
}

func (r *S3Registry) Initialize() error {
	s3client := s3.New(r.getSession())

	_, err := s3client.CreateBucket(&s3.CreateBucketInput{Bucket: aws.String(r.Bucket)})
	if err != nil {
		return err
	}
	return nil
}
