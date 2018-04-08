package registry

import (
	"errors"
	"io/ioutil"
	"log"
	"path"
	"os"
)

type FilesystemRegistry struct {
	BasePath string
}

func (r *FilesystemRegistry) ModuleExists(namespace,name,provider,version string) bool {
	return r.getModule(namespace,name,provider,version) != nil
}

func (r *FilesystemRegistry) GetModules() ([]Module, error) {
	namespaces, err := ioutil.ReadDir(r.BasePath)

	var result []Module
	if err != nil {
		return nil, err
	}

	for _, n := range namespaces {
		if !n.IsDir() {
			continue
		}

		modules, err := ioutil.ReadDir(path.Join(r.BasePath, n.Name()))

		if err != nil {
			return nil, err
		}

		for _, m := range modules {

			providers, err := ioutil.ReadDir(path.Join(r.BasePath, n.Name(), m.Name()))

			if err != nil {
				return nil, err
			}

			for _, p := range providers {

				versions, err := ioutil.ReadDir(path.Join(r.BasePath, n.Name(), m.Name(), p.Name()))

				if err != nil {
					return nil, err
				}

				for _, v := range versions {
					result = append(result, Module{
						ID:        path.Join(n.Name(), m.Name(), p.Name(), v.Name()),
						Name:      m.Name(),
						Namespace: n.Name(),
						Provider:  p.Name(),
						Version:   v.Name(),
					})
				}

			}
		}
	}

	return result, nil
}

func Filter(modules []Module, predicate func(Module) bool) []Module {
	var result []Module
	for _, module := range modules {
		if predicate(module) == true {
			result = append(result, module)
		}
	}

	return result
}

func (r *FilesystemRegistry) ListModules(namespace, name, provider string, offset, limit int) ([]Module, error) {

	modules, err := r.GetModules()

	modules = Filter(modules, func(m Module) bool {
		return (namespace == "" || m.Namespace == namespace) && (provider == "" || m.Provider == provider) && (name == "" || m.Name == name)
	})

	if err != nil {
		log.Fatal(err)
	}

	pageEnd := offset + limit

	if pageEnd > len(modules) {
		pageEnd = len(modules)
	}

	return modules[offset:pageEnd], nil
}
func (r *FilesystemRegistry) ListVersions(namespace, name, provider string) ([]ModuleVersions, error) {
	modules, _ := r.ListModules(namespace, name, provider, 0, 9999)
	versions := []ModuleVersions{{Source: namespace + "/" + name + "/" + provider, Versions: modules}}
	return versions, nil
}

func (r *FilesystemRegistry) searchModules(query string) error {
	return errors.New("Not implemented")
}

func (r *FilesystemRegistry) listVersions(namespace, name, provider string) error {
	return errors.New("Not implemented")
}

func (r *FilesystemRegistry) getDownloadUrl(namespace, name, provider, version string) error {
	return errors.New("Not implemented")
}

func (r *FilesystemRegistry) getLatestVersion(namespace, name string) error {
	return errors.New("Not implemented")
}

func (r *FilesystemRegistry) getLatestVersionByProvider(namespace, name, provider string) error {
	return errors.New("Not implemented")
}

func (r *FilesystemRegistry) getModule(namespace, name, provider, version string) *Module {
	_,err := os.Stat(path.Join(r.BasePath,namespace,name,provider,version))
	if os.IsNotExist(err) {
		return nil
	} else {
		return &Module {
			ID:        path.Join(namespace, name, provider, version),
			Name:      name,
			Namespace: namespace,
			Provider:  provider,
			Version:   version,
		}
	}
}
