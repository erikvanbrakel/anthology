package registry

import (
	"path"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"io"
	"bytes"
	"archive/tar"
	"compress/gzip"
)

type FilesystemRegistry struct {
	BasePath string
}

func (r *FilesystemRegistry) ListVersions(namespace, name, provider string) ([]ModuleVersions, error) {

	versions,err := r.getModules(namespace, name, provider)

	if err != nil {
		return nil,err
	}

	result := ModuleVersions {
		Source: strings.Join([]string{ namespace, name, provider }, "/"),
		Versions: versions,
	}

	return []ModuleVersions { result }, nil
}

func (r *FilesystemRegistry) ListModules(namespace, name, provider string, offset, limit int) ([]Module,int, error) {

	modules,err := r.getModules(namespace, name, provider)

	count := len(modules)

	if err != nil {
		return nil,0,err
	}

	if count == 0 {
		return modules[0:0],0,nil
	}

	end := limit + offset
	if (end) > len(modules) {
		end = len(modules)-1
	}

	return modules[offset:end],len(modules),nil
}

func (r *FilesystemRegistry) GetModule(namespace, name, provider, version string) (*Module, error) {
	_,err := os.Stat(path.Join(r.BasePath,namespace,name,provider,version))
	if err != nil {
		if os.IsNotExist(err) {
			return nil,nil
		} else {
			return nil,err
		}
	}
	return &Module {
		ID:        path.Join(namespace, name, provider, version),
		Name:      name,
		Namespace: namespace,
		Provider:  provider,
		Version:   version,
	}, nil
}

func (r *FilesystemRegistry) GetModuleData(namespace, name, provider, version string) (*bytes.Buffer, error) {

	module,_ := r.GetModule(namespace,name,provider,version)
	if module == nil {
		return nil,nil
	}

	tar, err := MakeTar(path.Join(r.BasePath, namespace, name, provider, version))
	gz := gzipit(tar)


	return  gz, err
}

func gzipit(source *bytes.Buffer) *bytes.Buffer {

	filename := "module.tar"
	var buffer bytes.Buffer

	archiver := gzip.NewWriter(&buffer)
	archiver.Name = filename
	defer archiver.Close()

	io.Copy(archiver, source)

	archiver.Flush()
	return &buffer
}

func MakeTar(inputPath string) (*bytes.Buffer, error) {

	var buffer bytes.Buffer

	tarball := tar.NewWriter(&buffer)
	defer tarball.Close()

	info, _ := os.Stat(inputPath)

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(inputPath)
	}

	e := filepath.Walk(inputPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name, _ = filepath.Rel(inputPath, path)
				if header.Name == "." {
					return nil
				}
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)

			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarball, file)
			return err
		},
	)
	tarball.Flush()
	return &buffer, e
}


func (r *FilesystemRegistry) getModules(namespace, name, provider string) ([]Module, error) {

	glob := r.BasePath

	if namespace != "" {
		glob = path.Join(glob,namespace)
	} else {
		glob = path.Join(glob,"*")
	}

	if name != "" {
		glob = path.Join(glob, name)
	} else {
		glob = path.Join(glob,"*")
	}

	if provider != "" {
		glob = path.Join(glob, provider)
	} else {
		glob = path.Join(glob,"*")
	}

	glob = path.Join(glob,"*")


	var modules []Module

	dirs,err := filepath.Glob(glob)

	if err != nil {
		return nil, errors.New("unable to read module directories")
	}

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

