package registry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/erikvanbrakel/anthology/app"
	"github.com/erikvanbrakel/anthology/models"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

type ArtifactoryRegistry struct {
	moduleURL   string
	providerURL string
}

type CheckSums struct {
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`
	SHA256 string `json:"sha256"`
}

type Child struct {
	URI    string `json:"uri"`
	Folder bool   `json:"folder"`
}

type Artifact struct {
	Repo        string    `json:"repo"`
	Path        string    `json:"path"`
	DownloadURI string    `json:"downloadUri"`
	Size        string    `json:"size"`
	CheckSums   CheckSums `json:"checksums"`
	Children    []Child   `json:"children"`
	URI         string    `json:"uri"`
}

func (r *ArtifactoryRegistry) ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error) {

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

func (r *ArtifactoryRegistry) PublishModule(namespace, name, provider, version string, data io.Reader) (err error) {
	panic("implement me")
}

func (r *ArtifactoryRegistry) GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error) {
	panic("implement me")
}

func (r *ArtifactoryRegistry) ListProviders(namespace, name string, offset, limit int) (providers []models.Provider, total int, err error) {

	providers, err = r.getProviders(namespace, name)

	count := len(providers)

	if err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return providers[0:0], 0, nil
	}

	end := limit + offset
	if (end) > len(providers) {
		end = len(providers)
	}

	return providers[offset:end], len(providers), nil

}

func (r *ArtifactoryRegistry) PublishProvider(namespace, name, version, OS, arch string, data io.Reader) (err error) {
	panic("implement me")
}

func (r *ArtifactoryRegistry) GetProviderData(namespace, name, version, OS, arch, file string) (reader *bytes.Buffer, err error) {

	// Get API information on the file
	artifactoryURL, err := url.Parse(r.providerURL)
	if err != nil {
		return nil, err
	}
	artifactoryURL.Path = path.Join(artifactoryURL.Path, namespace, name, version, file)

	a, err := getArtifact(artifactoryURL.String())
	artifactURL := a.DownloadURI

	// Get the file
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(artifactURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(respBody), nil
}

func (r *ArtifactoryRegistry) GetProviderMetaData(namespace, name, version, OS, arch string) (models.ProviderDownload, error) {

	artifactoryURL, err := url.Parse(r.providerURL)
	if err != nil {
		return models.ProviderDownload{}, err
	}
	artifactoryURL.Path = path.Join(artifactoryURL.Path, namespace, name, version)

	artifacts, err := walkArtifactory(artifactoryURL.String())

	if err != nil {
		return models.ProviderDownload{}, err
	}

	for _, a := range artifacts {
		urlPath, file := path.Split(a.Path)
		elements := strings.Split(urlPath, "/")
		if len(elements) < 3 {
			continue
		}
		fileParts := strings.Split(strings.TrimRight(file, ".zip"), "_")

		if len(fileParts) != 4 {
			continue
		}

		if OS != fileParts[2] || arch != fileParts[3] {
			continue
		}

		// Get SHASumsURL file
		SHASumsURL := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS", name, version)
		u, _ := url.Parse(r.providerURL)
		u.Path = path.Join(u.Path, namespace, name, version, SHASumsURL)
		_, err = getArtifact(u.String())
		if err != nil {
			// Log missing expected file
			logrus.Infof("Missing %s: %s", SHASumsURL, err)
			return models.ProviderDownload{}, err
		}

		// Get SHASumsSigURL file
		SHASumsSigURL := fmt.Sprintf("terraform-provider-%s_%s_SHA256SUMS.sig", name, version)
		u, _ = url.Parse(r.providerURL)
		u.Path = path.Join(u.Path, namespace, name, version, SHASumsSigURL)
		_, err = getArtifact(u.String())
		if err != nil {
			// Log missing expected file
			logrus.Infof("Missing %s: %s", SHASumsSigURL, err)
			return models.ProviderDownload{}, err
		}

		// Get GPG File
		gpgFile, err := r.GetProviderData(namespace, name, version, OS, arch, fmt.Sprintf("terraform-provider-%s_%s.gpg", name, version))
		gpgPublicKeys, err := getGPGKeys(gpgFile)
		if err != nil {
			logrus.Infof("Failed to load GPG file %s", err)
			return models.ProviderDownload{}, err
		}

		providerDownload := models.ProviderDownload{
			OS:            fileParts[2],
			Arch:          fileParts[3],
			Filename:      file,
			DownloadURL:   file,
			SHASumsURL:    SHASumsURL,
			SHASumsSigURL: SHASumsSigURL,
			SHASum:        a.CheckSums.SHA256,
			SigningKeys: struct {
				GPGPublicKeys []models.GPGKeys `json:"gpg_public_keys"`
			}{GPGPublicKeys: gpgPublicKeys},
		}

		return providerDownload, nil
	}

	return models.ProviderDownload{}, errors.New("Not Found")
}

func NewArtifactoryRegistry(options app.ArtifactoryOptions) Registry {

	registry := ArtifactoryRegistry{moduleURL: options.ModuleURL, providerURL: options.ProviderURL}

	if !strings.HasSuffix(registry.moduleURL, string(os.PathSeparator)) {
		registry.moduleURL = registry.moduleURL + string(os.PathSeparator)
	}
	if !strings.HasSuffix(registry.providerURL, string(os.PathSeparator)) {
		registry.providerURL = registry.providerURL + string(os.PathSeparator)
	}

	logrus.Infof("Using Artifactory Registry with ModuleURL %s and ProviderURL %s", registry.moduleURL, registry.providerURL)

	return &registry
}

func (r *ArtifactoryRegistry) getModules(namespace, name, provider string) ([]models.Module, error) {

	artifactoryURL, err := url.Parse(r.moduleURL)
	if err != nil {
		return nil, err
	}
	artifactoryURL.Path = path.Join(artifactoryURL.Path, namespace, name)

	artifacts, err := walkArtifactory(artifactoryURL.String())

	if err != nil {
		return nil, err
	}

	var modules []models.Module

	for _, a := range artifacts {
		path, file := path.Split(a.Path)
		elements := strings.Split(path, "/")
		version := elements[len(elements)-2]
		fileParts := strings.Split(strings.TrimRight(file, ".zip"), "_")

		if len(fileParts) != 4 {
			continue
		}

		modules = append(modules, models.Module{
			Namespace: namespace,
			Name:      name,
			Version:   version,
		})

	}

	return modules, nil
}

func (r *ArtifactoryRegistry) getProviders(namespace, name string) ([]models.Provider, error) {

	artifactoryURL, err := url.Parse(r.providerURL)
	if err != nil {
		return nil, err
	}
	artifactoryURL.Path = path.Join(artifactoryURL.Path, namespace, name)

	artifacts, err := walkArtifactory(artifactoryURL.String())

	if err != nil {
		return nil, err
	}

	var providers []models.Provider

	for _, a := range artifacts {
		path, file := path.Split(a.Path)
		elements := strings.Split(path, "/")
		if len(elements) < 3 {
			continue
		}
		if namespace == "" {
			namespace = elements[len(elements)-4]
		}
		if name == "" {
			name = elements[len(elements)-3]
		}
		version := elements[len(elements)-2]
		fileParts := strings.Split(strings.TrimRight(file, ".zip"), "_")

		if len(fileParts) != 4 {
			continue
		}

		providers = append(providers, models.Provider{
			Namespace: namespace,
			Name:      name,
			Version:   version,
			OS:        fileParts[2],
			Arch:      fileParts[3],
		})

	}

	return providers, nil
}

func getArtifact(url string) (Artifact, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	artifact := Artifact{}

	r, err := client.Get(url)
	if err != nil {
		return Artifact{}, err
	}
	defer r.Body.Close()

	respBody, err := ioutil.ReadAll(r.Body)

	err = json.Unmarshal(respBody, &artifact)

	return artifact, err
}

func walkArtifactory(url string) ([]Artifact, error) {
	artifacts := []Artifact{}

	a, err := getArtifact(url)

	if err != nil {
		return artifacts, err
	}

	if len(a.Children) > 0 {
		for _, c := range a.Children {
			art, err := walkArtifactory(a.URI + c.URI)
			if err != nil {

			} else {
				artifacts = append(artifacts, art...)
			}
		}
	} else {
		artifacts = append(artifacts, a)
	}

	return artifacts, nil
}
