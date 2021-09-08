package services

import (
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/erikvanbrakel/anthology/registry"
	"io"
)

type ProviderService struct {
	Registry registry.Registry
}

func NewProviderService(r registry.Registry) *ProviderService {
	return &ProviderService{
		r,
	}
}

func findProviderVersion(providers []models.ProviderVersions, namespace, name, version string) (int, bool) {
	for i, p := range providers {
		if namespace == p.Namespace && name == p.Name && version == p.Version {
			return i, true
		}
	}

	return -1, false
}

func (s *ProviderService) Query(rs app.RequestScope, namespace, name string, offset, limit int) ([]models.Provider, int, error) {

	providers, count, err := s.Registry.ListProviders(namespace, name, offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return providers, count, nil
}

func (s *ProviderService) QueryVersions(rs app.RequestScope, namespace, name string) ([]models.ProviderVersions, error) {

	providers, _, err := s.Registry.ListProviders(namespace, name, 0, 10000)

	if err != nil {
		return nil, err
	}

	// Convert to models.ProviderVersions
	var providerVersions []models.ProviderVersions
	for _, p := range providers {
		k, found := findProviderVersion(providerVersions, p.Namespace, p.Name, p.Version)
		if found {
			providerVersions[k].Platforms = append(providerVersions[k].Platforms, models.Platforms{p.OS, p.Arch})
		} else {
			providerVersions = append(providerVersions, models.ProviderVersions{
				p.Namespace,
				p.Name,
				p.Version,
				[]models.Platforms{models.Platforms{p.OS, p.Arch}},
			})
		}
	}

	return providerVersions, nil
}

func (s *ProviderService) Exists(rs app.RequestScope, namespace, name, version, OS, arch string) (bool, error) {
	providers, _, err := s.Registry.ListProviders(namespace, name, 0, 10000)

	if err != nil {
		return false, err
	}

	for _, p := range providers {
		if p.Version == version && p.OS == OS && p.Arch == arch {
			return true, nil
		}
	}

	return false, nil
}

func (s *ProviderService) Get(rs app.RequestScope, namespace, name, version, OS, arch string) (*models.Provider, error) {
	providers, _, err := s.Registry.ListProviders(namespace, name, 0, 10000)

	if err != nil {
		return nil, err
	}

	for _, p := range providers {
		if p.Version == version && p.OS == OS && p.Arch == arch {
			return &p, nil
		}
	}

	return nil, nil
}

func (s *ProviderService) Publish(rs app.RequestScope, namespace, name, version, OS, arch string, data io.Reader) error {
	return s.Registry.PublishProvider(namespace, name, version, OS, arch, data)
}

func (s *ProviderService) GetMetaData(rs app.RequestScope, namespace, name, version, OS, arch string) (models.ProviderDownload, error) {
	return s.Registry.GetProviderMetaData(namespace, name, version, OS, arch)
}

func (s *ProviderService) GetData(rs app.RequestScope, namespace, name, version, OS, arch, file string) (io.Reader, error) {
	return s.Registry.GetProviderData(namespace, name, version, OS, arch, file)
}
