package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FilesystemRegistry struct {
	BasePath string
}

func (r *FilesystemRegistry) getBasepath() string {

	if !strings.HasSuffix(r.BasePath, string(os.PathSeparator)) {
		return r.BasePath + string(os.PathSeparator)
	}
	return r.BasePath
}

func (r *FilesystemRegistry) ListVersions(namespace, name, provider string) ([]ModuleVersions, error) {

	versions, err := r.getModules(namespace, name, provider)

	if err != nil {
		return nil, err
	}

	result := ModuleVersions{
		Source:   strings.Join([]string{namespace, name, provider}, "/"),
		Versions: versions,
	}

	return []ModuleVersions{result}, nil
}

func (r *FilesystemRegistry) ListModules(namespace, name, provider string, offset, limit int) ([]Module, int, error) {

	modules, err := r.getModules(namespace, name, provider)

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

func (r *FilesystemRegistry) GetModule(namespace, name, provider, version string) (*Module, error) {
	_, err := os.Stat(path.Join(r.getBasepath(), namespace, name, provider, version) + ".tgz")

	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return &Module{
		ID:        path.Join(namespace, name, provider, version),
		Name:      name,
		Namespace: namespace,
		Provider:  provider,
		Version:   version,
	}, nil
}

func (r *FilesystemRegistry) GetModuleData(namespace, name, provider, version string) (*bytes.Buffer, error) {

	module, _ := r.GetModule(namespace, name, provider, version)
	if module == nil {
		return nil, nil
	}

	buffer := &bytes.Buffer{}

	f, err := os.Open(path.Join(r.getBasepath(), namespace, name, provider, version) + ".tgz")
	defer f.Close()

	if err != nil {
		return buffer, err
	}

	_, err = io.Copy(buffer, f)

	return buffer, err
}

func (r *FilesystemRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) (err error) {
	os.MkdirAll(path.Join(r.getBasepath(), namespace, name, provider), os.ModePerm)
	outfile, err := os.Create(path.Join(r.getBasepath(), namespace, name, provider, version) + ".tgz")
	defer outfile.Close()

	if err != nil {
		return err
	}

	_, err = io.Copy(outfile, data)
	return err
}

func (r *FilesystemRegistry) getModules(namespace, name, provider string) ([]Module, error) {

	glob := r.getBasepath()

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

	var modules []Module

	dirs, err := filepath.Glob(glob)

	if err != nil {
		return nil, errors.New("unable to read module directories")
	}

	for _, f := range dirs {
		parts := strings.Split(strings.TrimPrefix(f, r.getBasepath()), string(os.PathSeparator))

		if len(parts) != 4 {
			continue
		}

		modules = append(modules, Module{
			ID:        fmt.Sprintf("%s/%s/%s/%s", parts[0], parts[1], parts[2], strings.TrimRight(parts[3], ".tgz)")),
			Namespace: parts[0],
			Name:      parts[1],
			Provider:  parts[2],
			Version:   strings.TrimRight(parts[3], ".tgz"),
		})
	}

	return modules, nil
}
