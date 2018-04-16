package registry

type Registry interface {
	ListModules(namespace, name, provider string, offset, limit int) ([]Module, error)
	ListVersions(namespace, name, provider string) ([]ModuleVersions, error)
	GetModule(namespace, name, provider, version string) (*Module,  error)
}
