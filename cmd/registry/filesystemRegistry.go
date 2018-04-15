package registry

import (
	"path"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemRegistry struct {
	BasePath string
}

func (r *FilesystemRegistry) ListVersions(namespace, name, provider string) ([]ModuleVersions, error) {

	versions,_ := r.getModules(namespace, name, provider)

	result := ModuleVersions {
		Source: strings.Join([]string{ namespace, name, provider }, "/"),
		Versions: versions,
	}

	return []ModuleVersions { result }, nil
}

func (r *FilesystemRegistry) ListModules(namespace, name, provider string, offset, limit int) ([]Module, error) {

	modules,_ := r.getModules(namespace, name, provider)

	if (limit+offset) > len(modules) {
		limit = len(modules)-offset
	}

	return modules[offset:limit], nil
}

func (r *FilesystemRegistry) GetModule(namespace, name, provider, version string) *Module {
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

func (r *FilesystemRegistry) getModules(namespace, name, provider string) ([]Module, error) {

	basepath := r.BasePath

	if namespace != "" {
		basepath = path.Join(basepath,namespace)
	} else {
		basepath = path.Join(basepath,"*")
	}

	if name != "" {
		basepath = path.Join(basepath, name)
	} else {
		basepath = path.Join(basepath,"*")
	}

	if provider != "" {
		basepath = path.Join(basepath, provider)
	} else {
		basepath = path.Join(basepath,"*")
	}

	basepath = path.Join(basepath,"*")


	var modules []Module

	dirs,_ := filepath.Glob(basepath)

	for _,f := range dirs {
		parts := strings.Split(strings.TrimPrefix(f, r.BasePath), string(os.PathSeparator))

		if len(parts) != 4 {
			continue
		}

		modules = append(modules, Module{
			ID: strings.Join(parts, "/"),
			Namespace: parts[0],
			Name: parts[1],
			Provider: parts[2],
			Version: parts[3],
		})
	}

	return modules, nil
}

