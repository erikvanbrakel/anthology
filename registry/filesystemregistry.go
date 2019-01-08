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
	"io/ioutil"
	"fmt"
)

type FilesystemRegistry struct {
	basePath string
}

func NewFilesystemRegistry(basePath string) (Registry, error) {
	registry := &FilesystemRegistry{basePath }

	if !strings.HasSuffix(registry.basePath, string(os.PathSeparator)) {
		registry.basePath = registry.basePath + string(os.PathSeparator)
	}

	return registry, nil
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

func (r *FilesystemRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) (error) {
	content, _ := ioutil.ReadAll(data)
	os.MkdirAll(fmt.Sprintf("%v/%v/%v/%v/", r.basePath, namespace, name, provider), os.ModePerm)
	return ioutil.WriteFile(fmt.Sprintf("%v/%v/%v/%v/%v.tgz", r.basePath, namespace, name, provider, version), content, os.ModePerm)
}

func (r *FilesystemRegistry) GetModuleData(namespace, name, provider, version string) (*bytes.Buffer, error) {

	content, _ := ioutil.ReadFile(fmt.Sprintf("%v/%v/%v/%v/%v.tgz", r.basePath, namespace, name, provider, version))

	return bytes.NewBuffer(content), nil

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
			Data:      func() (*bytes.Buffer, error) {
				return r.GetModuleData(namespace, name, provider, strings.TrimRight(parts[3], ".tgz"))
			},
			Version:   strings.TrimRight(parts[3], ".tgz"),
		})
	}

	return modules, nil
}
