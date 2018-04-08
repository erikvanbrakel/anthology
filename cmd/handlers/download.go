package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func DownloadHandler(basePath string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		namespace, name, provider, version := params["namespace"], params["name"], params["provider"], params["version"]

		tar, _ := MakeTar(path.Join(basePath, namespace, name, provider, version))

		gz := gzipit(tar)

		w.Header().Set("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s-%s-%s-%s.tgz", namespace, name, provider, version))

		written, _ := io.Copy(w, gz)

		w.Header().Set("Content-Length", strconv.FormatInt(written, 10))
	}
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
