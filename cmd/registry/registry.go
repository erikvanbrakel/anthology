package registry

type Registry interface {
	ListModules(namespace, name, provider string, offset, limit int) (modules []Module, total int, err error)
	ListVersions(namespace, name, provider string) (versions []ModuleVersions, err error)
	GetModule(namespace, name, provider, version string) (module *Module,  err error)
}
