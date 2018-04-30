package registry

import (
	"bytes"
	"github.com/erikvanbrakel/anthology/models"
	"io"
)

type Registry interface {
	/*	ListModules(namespace, name, provider string, offset, limit int) (modules []Module, total int, err error)
		ListVersions(namespace, name, provider string) (versions []ModuleVersions, err error)
		GetModule(namespace, name, provider, version string) (module *Module, err error) */
	GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error)
	ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error)
	PublishModule(namespace, name, provider, version string, data io.Reader) (err error)
}
