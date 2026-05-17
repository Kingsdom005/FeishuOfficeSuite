package registry

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
)

type Registry struct {
	consulAddress string
}

func NewRegistry(consulAddress string) *Registry {
	return &Registry{
		consulAddress: consulAddress,
	}
}

func (r *Registry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return nil
}

func (r *Registry) GetService(ctx context.Context, name string) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

func (r *Registry) ListServices(ctx context.Context) ([]*registry.ServiceInstance, error) {
	return nil, nil
}

type ServiceRegistry struct {
	registry.Registry
}

func NewServiceRegistry(r registry.Registry) *ServiceRegistry {
	return &ServiceRegistry{Registry: r}
}