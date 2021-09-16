package registry

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/sirupsen/logrus"
	"io"
	"path"
	"strings"
)

type S3Registry struct {
	bucket   string
	endpoint string
}

func (r *S3Registry) ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error) {
	modules, err = r.getModules(namespace, name, provider)

	if err != nil {
		return nil, 0, err
	}

	return modules, len(modules), nil
}

func (S3Registry) PublishModule(namepsace, name, provider, version string, data io.Reader) (err error) {
	panic("implement me")
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
				Version:   strings.TrimRight(parts[3], ".tgz"),
			})
		}
	}

	return modules, nil
}

func (r *S3Registry) ListProviders(namespace, name string, offset, limit int) (providers []models.Provider, total int, err error) {
	providers, err = r.getProviders(namespace, name)

	if err != nil {
		return nil, 0, err
	}

	return providers, len(providers), nil
}

func (S3Registry) PublishProvider(namepsace, name, version, OS, arch string, data io.Reader) (err error) {
	panic("implement me")
}

func (r *S3Registry) GetProviderData(namespace, name, version, OS, arch, file string) (reader *bytes.Buffer, err error) {
	s3client := s3.New(r.getSession())

	key := strings.Join([]string{namespace, name, version, file}, "/")
	obj, err := s3client.GetObject(&s3.GetObjectInput{
		Key:    aws.String(key),
		Bucket: aws.String(r.bucket),
	})

	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	io.Copy(buffer, obj.Body)
	return buffer, nil
}

func (r *S3Registry) GetProviderMetaData(namespace, name, version, OS, arch string) (providerMetaData models.ProviderDownload, err error) {
	s3client := s3.New(r.getSession())

	prefix := path.Join(namespace, name, version, "*")

	loi := s3.ListObjectsInput{
		Bucket: aws.String(r.bucket),
		Prefix: aws.String(prefix),
	}
	result, err := s3client.ListObjects(&loi)

	if err != nil {
		logrus.Errorf("error: %s", err.(awserr.Error))
		return models.ProviderDownload{}, err
	}

	providerDownload := models.ProviderDownload{
		OS:   OS,
		Arch: arch,
	}

	for _, o := range result.Contents {
		// <namespace>/<name>/<version>/terraform-provider-<name>_<version>_<OS>_<arch>.zip
		parts := strings.Split(*o.Key, "/")

		if len(parts) != 4 {
			continue
		}

		if parts[3] == fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, OS, arch) {
			providerDownload.Filename = parts[3]
			providerDownload.DownloadURL = parts[3]

			// Generate SHASum
			shaSUMFile, err := r.GetProviderData(namespace, name, version, OS, arch, parts[3])
			shasum, err := sha256File(shaSUMFile)
			if err != nil {
				logrus.Infof("Failed to SHA256 file")
				return models.ProviderDownload{}, err
			}
			providerDownload.SHASum = shasum

		} else if parts[3] == fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", name, version) {
			providerDownload.SHASumsURL = parts[3]

		} else if parts[3] == fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", name, version) {
			providerDownload.SHASumsSigURL = parts[3]

		} else if parts[3] == fmt.Sprintf("terraform-provider-%s_%s.gpg", name, version) {
			// GPG Public Keys
			gpgFile, err := r.GetProviderData(namespace, name, version, OS, arch, parts[3])
			gpgPublicKeys, err := getGPGKeys(gpgFile)
			if err != nil {
				logrus.Infof("Failed to load GPG file %s", err)
				return models.ProviderDownload{}, err
			}
			providerDownload.SigningKeys = struct {
				GPGPublicKeys []models.GPGKeys `json:"gpg_public_keys"`
			}{GPGPublicKeys: gpgPublicKeys}
		}

	}

	if providerDownload.Filename != "" && providerDownload.SHASumsURL != "" && providerDownload.SHASumsSigURL != "" {
		return providerDownload, nil
	}

	return models.ProviderDownload{}, errors.New("Not Found")

}

func (r *S3Registry) getProviders(namespace, name string) (providers []models.Provider, err error) {
	prefix := ""

	if namespace != "" {
		prefix = namespace
		if name != "" {
			prefix = strings.Join([]string{prefix, name}, "/")
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
		// <namespace>/<name>/<version>/terraform-provider-<name>_<version>_<OS>_<arch>.zip
		parts := strings.Split(*o.Key, "/")

		if len(parts) != 4 {
			continue
		}

		fileParts := strings.Split(strings.TrimRight(parts[3], ".zip"), "_")

		if len(fileParts) != 4 {
			continue
		}

		providers = append(providers, models.Provider{
			Namespace: parts[0],
			Name:      parts[1],
			Version:   parts[2],
			OS:        fileParts[2],
			Arch:      fileParts[3],
		})
	}

	return providers, nil
}

func (r *S3Registry) getSession() *session.Session {
	config := &aws.Config{
		S3ForcePathStyle: aws.Bool(true),
		Region:           aws.String("us-east-1"),
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

func NewS3Registry(options app.S3Options) Registry {
	return &S3Registry{
		options.Bucket,
		options.Endpoint,
	}
}
