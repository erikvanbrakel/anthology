package tests

import (
	"testing"

	ctls "crypto/tls"

    "gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/gentleman.v2/plugins/tls"
	"archive/tar"
	"compress/gzip"
	"bytes"
	"os"
	"io"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"path"
	"net/http"
	"os/exec"
	"time"
)

func ComposeUp() error {
	path, err := exec.LookPath("docker-compose")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "up", "-d")
	cmd.Start()
	err = cmd.Wait()

	if err != nil {
		return err
	}

	return nil
}

func ComposeDown() {
	cmd := exec.Command("docker-compose", "down")
	cmd.Start()
	cmd.Wait()
}

func TestMain(m *testing.M) {
	ComposeUp()
	time.Sleep(time.Second * 10)
	retCode := m.Run()
	ComposeDown()
	os.Exit(retCode)
}

func TestIntegration(t *testing.T) {

	registries := map[string]string {
		"s3": "https://registry-s3.anthology.localtest.me",
		"filesystem": "https://registry-filesystem.anthology.localtest.me",
	}

	for k,v := range registries {
		t.Run(k, func(t *testing.T) {
			var test = baloo.New(v)
			test.Use(tls.Config(&ctls.Config{InsecureSkipVerify: true}))

			t.Run("check disco file", func(t *testing.T) {
				test.Get("/.well-known/terraform.json").
					Expect(t).
					Type("json").
					Done()
			})

			t.Run("verify registry is empty", func (t *testing.T) {
				// initially empty
				test.Get("/v1/").
					Expect(t).
					Status(http.StatusNotFound).
					Done()
			})

			t.Run("upload module", func (t *testing.T) {
				// Add a module
				tarball, err := createModuleTarball()
				if err != nil {
					assert.Fail(t, err.Error())
				}
				test.Post("/v1/anthology.tests/mymodule/aws/1.0.0").
					Body(tarball).
					Expect(t).
					Status(http.StatusNoContent).
					Done()
			})

			t.Run("verify module exists", func(t *testing.T) {
				// now we can find things
				test.Get("/v1/anthology.tests/mymodule/aws/1.0.0").
					Expect(t).
					Status(http.StatusOK).
					Type("json").
					Done()
			})

		})
	}

}

func createModuleTarball() (io.Reader, error) {

	buffer := bytes.NewBuffer([]byte{})
	gzipWriter := gzip.NewWriter(buffer)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	files,_ := ioutil.ReadDir("../tests/terraform-module")

	for _,f := range files {
		fp := path.Join("../tests/terraform-module", f.Name())
		file, err := os.Open(fp)
		defer file.Close()

		if err != nil {
			return nil, err
		}

		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}

		header := &tar.Header{
			Name:    fp,
			Size:    stat.Size(),
			Mode:    int64(stat.Mode()),
			ModTime: stat.ModTime(),
		}

		err = tarWriter.WriteHeader(header)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return nil, err
		}

	}
	return buffer, nil
}