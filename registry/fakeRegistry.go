package registry

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/erikvanbrakel/anthology/models"
	"io"
	"strings"
)

type InMemoryRegistry struct {
	modules   []models.Module
	providers []models.Provider
	data      map[string][]byte
}

func (r *InMemoryRegistry) ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error) {
	var result []models.Module

	for _, m := range r.modules {
		if namespace != "" && m.Namespace != namespace {
			continue
		}
		if name != "" && m.Name != name {
			continue
		}
		if provider != "" && m.Provider != provider {
			continue
		}
		result = append(result, m)
	}
	return result, len(result), nil
}

func (r *InMemoryRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) error {
	r.modules = append(r.modules, models.Module{
		namespace,
		name,
		provider,
		version,
	})

	id := strings.Join([]string{namespace, name, provider, version}, "/")

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)

	r.data[id] = buf.Bytes()

	return nil
}

func (r *InMemoryRegistry) GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error) {
	id := strings.Join([]string{namespace, name, provider, version}, "/")
	moduleData, exists := r.data[id]
	if !exists {
		return nil, errors.New("module does not exist")
	}

	return bytes.NewBuffer(moduleData), nil
}

func (r *InMemoryRegistry) ListProviders(namespace, name string, offset, limit int) (providers []models.Provider, total int, err error) {
	var result []models.Provider

	for _, p := range r.providers {
		if namespace != "" && p.Namespace != namespace {
			continue
		}
		if name != "" && p.Name != name {
			continue
		}
		result = append(result, p)
	}
	return result, len(result), nil
}

func (r *InMemoryRegistry) PublishProvider(namespace, name, version, OS, arch string, data io.Reader) error {
	r.providers = append(r.providers, models.Provider{
		namespace,
		name,
		version,
		OS,
		arch,
	})

	id := strings.Join([]string{namespace, name, version, OS, arch}, "/")

	buf := new(bytes.Buffer)
	buf.ReadFrom(data)

	r.data[id] = buf.Bytes()

	return nil
}
func (r *InMemoryRegistry) GetProviderMetaData(namespace, name, version, OS, arch string) (providerMetaData models.ProviderDownload, err error) {
	id := strings.Join([]string{namespace, name, version, OS, arch}, "/")
	providerData, exists := r.data[id]
	if !exists {
		return models.ProviderDownload{}, errors.New("provider does not exist")
	}

	providerDownload := models.ProviderDownload{
		OS:            OS,
		Arch:          arch,
		Filename:      fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, OS, arch),
		DownloadURL:   "",
		SHASumsURL:    "",
		SHASumsSigURL: "",
		SHASum:        fmt.Sprintf("%x", sha256.Sum256(providerData)),
	}

	return providerDownload, nil
}

func (r *InMemoryRegistry) GetProviderData(namespace, name, version, OS, arch, file string) (reader *bytes.Buffer, err error) {
	id := strings.Join([]string{namespace, name, version, OS, arch}, "/")
	providerData, exists := r.data[id]
	if !exists {
		return nil, errors.New("provider does not exist")
	}
	if file == fmt.Sprintf("terraform-provider-%s_%s_%s_%s.zip", name, version, OS, arch) {
		return bytes.NewBuffer(providerData), nil
	} else if file == fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", name, version) {
		fileContents := fmt.Sprintf("%x  terraform-provider-%s_%s_%s_%s.zip", sha256.Sum256(providerData), name, version, OS, arch)
		return bytes.NewBuffer([]byte(fileContents)), nil
	}

	return nil, errors.New("requested file does not exist")
}

func NewFakeRegistry() Registry {
	return &InMemoryRegistry{data: map[string][]byte{}}
}
