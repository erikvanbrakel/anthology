package registry

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"errors"
)

type S3Registry struct {
	bucket   string
	endpoint string
	region   string
}

func NewS3Registry(bucket, endpoint, region string) (Registry, error) {
	if bucket == "" { return nil, errors.New("bucket doesn't exist") }

	return &S3Registry{
		bucket: bucket,
		endpoint: endpoint,
		region: region,
	}, nil
}

func (r *S3Registry) ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error) {
	modules, err = r.getModules(namespace, name, provider)

	if err != nil {
		return nil, 0, err
	}

	return modules, len(modules), nil
}

func (r *S3Registry) PublishModule(namespace, name, provider, version string, data io.Reader) (error) {
	s3client := s3.New(r.getSession())
	manager := s3manager.NewUploaderWithClient(s3client)


	_, err := manager.Upload(&s3manager.UploadInput {
		Body:   data,
		Key:    aws.String(strings.Join([]string{namespace, name, provider, version}, "/") + ".tgz"),
		Bucket: aws.String(r.bucket),
	})

	return err
}

func (r *S3Registry) GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error) {
	s3client := s3.New(r.getSession())

	obj, err := s3client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(strings.Join([]string{namespace, name, provider, version}, "/") + ".tgz"),
		Bucket: aws.String(r.bucket),
	})

	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	io.Copy(buffer, obj.Body)
	return buffer, nil
}

func (r *S3Registry) getModules(namespace, name, provider string) (modules []models.Module, err error) {
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

	if prefix != "" {
		prefix += "/"
	}

	s3client := s3.New(r.getSession())

	loi := s3.ListObjectsInput{
		Bucket: aws.String(r.bucket),
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
			modules = append(modules, models.Module{
				Namespace: parts[0],
				Name:      parts[1],
				Provider:  parts[2],
				Data:      func() (*bytes.Buffer, error) {
					return r.GetModuleData(namespace, name, provider, strings.TrimRight(parts[3], ".tgz"))
				},
				Version:   strings.TrimRight(parts[3], ".tgz"),
			})
		}
	}

	return modules, nil
}

func (r *S3Registry) getSession() *session.Session {
	config := &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String(r.region),
	}
	if r.endpoint != "" {
		if !strings.HasPrefix(r.endpoint, "https") {
			config.DisableSSL = aws.Bool(true)
		}
		config.Endpoint = aws.String(r.endpoint)
	}
	config.CredentialsChainVerboseErrors = aws.Bool(true)

	s, _ := session.NewSession(config)

	return s
}

func (r *S3Registry) Initialize() error {


	s3client := s3.New(r.getSession())

	_, err := s3client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(r.bucket),
	})

	return err
}