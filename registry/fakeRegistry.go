package registry

import (
	"bytes"
	"errors"
	"github.com/erikvanbrakel/anthology/models"
	"io"
	"strings"
)

type InMemoryRegistry struct {
	modules []models.Module
	data    map[string][]byte
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

	contents := new(bytes.Buffer)
	contents.ReadFrom(data)

	r.modules = append(r.modules, models.Module{
		Namespace: namespace,
		Name:      name,
		Provider:  provider,
		Version:   version,
		Data:     func() (*bytes.Buffer,error) {
			return contents, nil
		},
	})

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

func NewFakeRegistry() Registry {
	return &InMemoryRegistry{data: map[string][]byte{}}
}

func (r *InMemoryRegistry) Initialize() error {
	return nil
}