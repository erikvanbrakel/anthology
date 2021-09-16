package registry

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FilesystemRegistry struct {
	basePath     string
	providerPath string
}

func (r *FilesystemRegistry) ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error) {

	modules, err = r.getModules(namespace, name, provider)

	count := len(modules)

	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return modules[0:0], 0, nil
	}

	end := limit + offset
	if (end) > len(modules) {
		end = len(modules)
	}

	return modules[offset:end], len(modules), nil

}

func (r *FilesystemRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) (err error) {
	panic("implement me")
}

func (r *FilesystemRegistry) GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error) {
	panic("implement me")
}

func (r *FilesystemRegistry) ListProviders(namespace, name string, offset, limit int) (providers []models.Provider, total int, err error) {

	providers, err = r.getProviders(namespace, name)

	count := len(providers)

	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return providers[0:0], 0, nil
	}

	end := limit + offset
	if (end) > len(providers) {
		end = len(providers)
	}

	return providers[offset:end], len(providers), nil

}

func (r *FilesystemRegistry) PublishProvider(namespace, name, version, OS, arch string, data io.Reader) (err error) {
	panic("implement me")
}

func (r *FilesystemRegistry) GetProviderData(namespace, name, version, OS, arch, file string) (reader *bytes.Buffer, err error) {
	filePath := path.Join(r.providerPath, namespace, name, version, file)

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader = new(bytes.Buffer)
	reader.ReadFrom(f)

	return reader, nil
}

func (r *FilesystemRegistry) GetProviderMetaData(namespace, name, version, OS, arch string) (providerMetaData models.ProviderDownload, err error) {

	// Verify zip file exists
	zipFile := fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, OS, arch)
	if _, err := os.Stat(path.Join(r.providerPath, namespace, name, version, zipFile)); os.IsNotExist(err) {
		return models.ProviderDownload{}, err
	}

	// Verify SHASums file exists
	SHASumsURL := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", name, version)
	if _, err := os.Stat(path.Join(r.providerPath, namespace, name, version, SHASumsURL)); os.IsNotExist(err) {
		// Log missing expected file
		logrus.Infof("Missing %s: %s", SHASumsURL, err)
		return models.ProviderDownload{}, err
	}

	// Verify SHASumsSig file exists
	SHASumsSigURL := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", name, version)
	if _, err := os.Stat(path.Join(r.providerPath, namespace, name, version, SHASumsSigURL)); os.IsNotExist(err) {
		// Log missing expected file
		logrus.Infof("Missing %s: %s", SHASumsSigURL, err)
		return models.ProviderDownload{}, err
	}

	// Verify GPG file exists
	gpgFile := path.Join(r.providerPath, namespace, name, version, fmt.Sprintf("terraform-provider-%s_%s.gpg", name, version))
	if _, err := os.Stat(gpgFile); os.IsNotExist(err) {
		// Log missing expected file
		logrus.Infof("Missing %s: %s", gpgFile, err)
		return models.ProviderDownload{}, err
	}
	gpgFH, _ := os.Open(gpgFile)
	defer gpgFH.Close()
	gpgPublicKeys, err := getGPGKeys(gpgFH)

	if err != nil {
		logrus.Infof("Failed to load GPG file %s", err)
		return models.ProviderDownload{}, err
	}

	fileFH, _ := os.Open(gpgFile)
	defer fileFH.Close()
	shasum, err := sha256File(fileFH)
	if err != nil {
		logrus.Infof("Failed to SHA256 file")
		return models.ProviderDownload{}, err
	}

	providerDownload := models.ProviderDownload{
		OS:            OS,
		Arch:          arch,
		Filename:      zipFile,
		DownloadURL:   zipFile,
		SHASumsURL:    SHASumsURL,
		SHASumsSigURL: SHASumsSigURL,
		SHASum:        shasum,
		SigningKeys: struct {
			GPGPublicKeys []models.GPGKeys `json:"gpg_public_keys"`
		}{
			GPGPublicKeys: gpgPublicKeys},
	}

	return providerDownload, nil

}

func NewFilesystemRegistry(options app.FileSystemOptions) Registry {

	registry := FilesystemRegistry{basePath: options.BasePath, providerPath: options.ProviderPath}

	if !strings.HasSuffix(registry.basePath, string(os.PathSeparator)) {
		registry.basePath = registry.basePath + string(os.PathSeparator)
	}
	if !strings.HasSuffix(registry.providerPath, string(os.PathSeparator)) {
		registry.providerPath = registry.providerPath + string(os.PathSeparator)
	}

	logrus.Infof("Using Filesystem Registry with basepath %s and providerpath %s", registry.basePath, registry.providerPath)

	return &registry
}

func (r *FilesystemRegistry) getModules(namespace, name, provider string) ([]models.Module, error) {

	glob := r.basePath

	if namespace != "" {
		glob = path.Join(glob, namespace)
	} else {
		glob = path.Join(glob, "*")
	}

	if name != "" {
		glob = path.Join(glob, name)
	} else {
		glob = path.Join(glob, "*")
	}

	if provider != "" {
		glob = path.Join(glob, provider)
	} else {
		glob = path.Join(glob, "*")
	}

	glob = path.Join(glob, "*.tgz")

	var modules []models.Module

	dirs, err := filepath.Glob(glob)

	if err != nil {
		return nil, errors.New("unable to read module directories")
	}

	for _, f := range dirs {
		parts := strings.Split(strings.TrimPrefix(f, r.basePath), string(os.PathSeparator))

		if len(parts) != 4 {
			continue
		}

		modules = append(modules, models.Module{
			Namespace: parts[0],
			Name:      parts[1],
			Provider:  parts[2],
			Version:   strings.TrimRight(parts[3], ".tgz"),
		})
	}

	return modules, nil
}

func (r *FilesystemRegistry) getProviders(namespace, name string) ([]models.Provider, error) {

	glob := r.providerPath

	if namespace != "" {
		glob = path.Join(glob, namespace)
	} else {
		glob = path.Join(glob, "*")
	}

	if name != "" {
		glob = path.Join(glob, name)
	} else {
		glob = path.Join(glob, "*")
	}

	glob = path.Join(glob, "*/*.zip")

	var providers []models.Provider

	dirs, err := filepath.Glob(glob)

	if err != nil {
		return nil, errors.New("unable to read provider directories")
	}

	for _, f := range dirs {
		parts := strings.Split(strings.TrimPrefix(f, r.providerPath), string(os.PathSeparator))

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
