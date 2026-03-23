package cli

type ModelCLIService interface {
	Name() string
	Service() any
}

type ServiceRegistry struct {
	services map[string]any
}

var (
	cli_services []ModelCLIService
)

func NewServiceRegistry(modules []ModelCLIService) *ServiceRegistry {
	m := make(map[string]any)

	for _, module := range modules {
		m[module.Name()] = module.Service()
	}

	return &ServiceRegistry{
		services: m,
	}
}
