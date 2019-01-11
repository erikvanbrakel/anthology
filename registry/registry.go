package registry

import (
	"bytes"
	"github.com/erikvanbrakel/anthology/models"
	"io"
)

type Registry interface {
	Initialize() error
	GetModuleData(namespace, name, provider, version string) (reader *bytes.Buffer, err error)
	ListModules(namespace, name, provider string, offset, limit int) (modules []models.Module, total int, err error)
	PublishModule(namespace, name, provider, version string, data io.Reader) (err error)
}
