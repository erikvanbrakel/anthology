package registry

import (
	"bytes"
	"errors"
	"github.com/erikvanbrakel/anthology/models"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FilesystemRegistry struct {
	basePath string
}

func NewFilesystemRegistry(basePath string) (Registry, error) {
	return &FilesystemRegistry{
		basePath: basePath,
	}, nil
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

/*func NewFilesystemRegistry(options app.FileSystemOptions) Registry {

	registry := FilesystemRegistry{basePath: options.BasePath}

	if !strings.HasSuffix(registry.basePath, string(os.PathSeparator)) {
		registry.basePath = registry.basePath + string(os.PathSeparator)
	}

	logrus.Infof("Using Filesystem Registry with basepath %s", registry.basePath)

	return &registry
}*/

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
