package services

import (
	"github.com/erikvanbrakel/anthology/models"
	"github.com/erikvanbrakel/anthology/registry"
	"io"
	"errors"
)

type ModuleService interface {
	Query(namespace, name, provider string, verified bool, offset, limit int) ([]models.Module, int, error)
	QueryVersions(namespace, name, provider string) ([]models.Module, error)
	Exists(namespace, name, provider, version string) (bool, error)
	Get(namespace, name, provider, version string) (*models.Module, error)
	GetData(namespace, name, provider, version string) (io.Reader, error)
	Publish(namespace, name, provider, version string, data io.Reader) error
}


type moduleService struct {
	Registry registry.Registry
}

func NewModuleService(r registry.Registry) (*moduleService, error) {
	if r == nil {
		return nil, errors.New("registry can't be nil")
	}

	return &moduleService{
		r,
	}, nil
}

func (s *moduleService) Query(namespace, name, provider string, verified bool, offset, limit int) ([]models.Module, int, error) {

	modules, count, err := s.Registry.ListModules(namespace, name, provider, offset, limit)

	if err != nil {
		return nil, 0, err
	}

	return modules, count, nil
}

func (s *moduleService) QueryVersions(namespace, name, provider string) ([]models.Module, error) {

	modules, _, err := s.Registry.ListModules(namespace, name, provider, 0, 10000)
	return modules, err
}

func (s *moduleService) Exists(namespace, name, provider, version string) (bool, error) {
	modules, _, err := s.Registry.ListModules(namespace, name, provider, 0, 10000)

	if err != nil {
		return false, err
	}

	for _, m := range modules {
		if m.Version == version {
			return true, nil
		}
	}
	return false, nil
}

func (s *moduleService) Get(namespace, name, provider, version string) (*models.Module, error) {
	modules, _, err := s.Registry.ListModules(namespace, name, provider, 0, 10000)

	if err != nil {
		return nil, err
	}

	for _, m := range modules {
		if m.Version == version {
			return &m, nil
		}
	}

	return nil, nil
}

func (s *moduleService) Publish(namespace, name, provider, version string, data io.Reader) error {
	return s.Registry.PublishModule(namespace, name, provider, version, data)
}

func (s *moduleService) GetData(namespace, name, provider, version string) (io.Reader, error) {
	return s.Registry.GetModuleData(namespace, name, provider, version)
}
