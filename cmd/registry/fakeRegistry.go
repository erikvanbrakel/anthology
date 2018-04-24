package registry

import (
	"bytes"
	"io"
	"strings"
)

type FakeRegistry struct {
	Modules []Module
}

func (r *FakeRegistry) ListModules(namespace, name, provider string, offset, limit int) (modules []Module, total int, err error) {

	for _, m := range r.Modules {
		if namespace != "" && m.Namespace != namespace {
			continue
		}
		if name != "" && m.Name != name {
			continue
		}
		if provider != "" && m.Provider != provider {
			continue
		}

		modules = append(modules, m)
	}

	return modules, len(modules), nil
}

func (r *FakeRegistry) ListVersions(namespace, name, provider string) (versions []ModuleVersions, err error) {
	m, _, _ := r.ListModules(namespace, name, provider, 0, len(r.Modules))
	return []ModuleVersions{
		{
			Source:   strings.Join([]string{namespace, name, provider}, "/"),
			Versions: m,
		},
	}, nil
}

func (r *FakeRegistry) GetModule(namespace, name, provider, version string) (module *Module, err error) {
	for _, m := range r.Modules {
		if namespace != "" && m.Namespace != namespace {
			continue
		}
		if name != "" && m.Name != name {
			continue
		}
		if provider != "" && m.Provider != provider {
			continue
		}
		if version != "" && m.Version != version {
			continue
		}

		return &m, nil
	}

	return nil, nil
}

func (r *FakeRegistry) GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error) {
	return &bytes.Buffer{}, nil
}

func (r *FakeRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) (err error) {
	r.Modules = append(r.Modules, Module{
		ID:        strings.Join([]string{namespace, name, provider, version}, "/"),
		Namespace: namespace,
		Name:      name,
		Provider:  provider,
		Version:   version,
	})
	return nil
}
