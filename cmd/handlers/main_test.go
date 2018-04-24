package handlers_test

import (
	"fmt"
	"github.com/erikvanbrakel/anthology/cmd"
	"github.com/erikvanbrakel/anthology/cmd/registry"
	"github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var servers []struct {
	Name   string
	Server httptest.Server
}

func TestMain(m *testing.M) {

	access_key, secret_key := "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"

	os.Setenv("AWS_ACCESS_KEY", access_key)
	os.Setenv("AWS_SECRET_KEY", secret_key)

	dispose, endpoint := createMinioContainer(access_key, secret_key)
	defer dispose()

	s3r := &registry.S3Registry{Bucket: "modules", Endpoint: endpoint}
	s3r.Initialize()
	fakeserver, _ := cmd.NewServer(cmd.RegistryServerConfig{}, &registry.FakeRegistry{})
	s3server, _ := cmd.NewServer(cmd.RegistryServerConfig{}, s3r)
	fs_basedir, _ := ioutil.TempDir("", "fs")
	filesystemserver, _ := cmd.NewServer(cmd.RegistryServerConfig{}, &registry.FilesystemRegistry{BasePath: fs_basedir})

	servers = []struct {
		Name   string
		Server httptest.Server
	}{
		{
			Name:   "FakeRegistry",
			Server: *httptest.NewServer(fakeserver.Router),
		},
		{
			Name:   "S3",
			Server: *httptest.NewServer(s3server.Router),
		},
		{
			Name:   "FileSystem",
			Server: *httptest.NewServer(filesystemserver.Router),
		},
	}

	m.Run()
}

func createMinioContainer(access_key, secret_key string) (dispose func(), endpoint string) {

	client, err := docker.NewClientFromEnv()
	if err != nil {
		os.Exit(1)
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		HostConfig: &docker.HostConfig{
			PublishAllPorts: true,
		},
		Config: &docker.Config{
			Image: "minio/minio",
			Cmd:   []string{"server", "/data"},
			Env: []string{
				"MINIO_ACCESS_KEY=" + access_key,
				"MINIO_SECRET_KEY=" + secret_key,
			},
		},
	})

	if err != nil {
		os.Exit(1)
	}

	err = client.StartContainer(c.ID, &docker.HostConfig{})
	if err != nil {
		os.Exit(1)
	}

	dispose = func() {
		err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    c.ID,
			Force: true,
		})

		if err != nil {
			os.Exit(1)
		}
	}

	for i := 0; i < 30; i++ {
		s, _ := client.InspectContainer(c.ID)
		if s.State.Health.Status == "healthy" {
			binding := s.NetworkSettings.Ports["9000/tcp"][0]
			return dispose, fmt.Sprintf("http://localhost:%s", binding.HostPort)
		}
		logrus.Infof("Current status: %s", s.State.Health.Status)
		time.Sleep(time.Second * 5)
	}

	dispose()
	os.Exit(1)

	return nil, ""
}
