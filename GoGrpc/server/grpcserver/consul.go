package grpcserver

import (
	"fmt"
	"server/config"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type ConsulServer struct {
	cfg    *config.ConsulConfig
	client *api.Client
}

func NewConsulServer(c *config.ConsulConfig) (*ConsulServer, error) {
	if c == nil {
		return nil, errors.New("consul config is nil")
	}
	consulCfg := api.DefaultConfig()
	if c.Scheme != "" {
		consulCfg.Scheme = c.Scheme
	}
	if c.Addr != "" {
		consulCfg.Address = c.Addr
	}
	client, err := api.NewClient(consulCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %s\n", err)
	}
	return &ConsulServer{cfg: c, client: client}, nil
}

func (cs *ConsulServer) RegisterService(serviceName, ip string, port int, tags []string) error {
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, port),
		Name:    serviceName,
		Tags:    tags,
		Address: ip,
		Port:    port,
	}
	if err := cs.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service to consul: %s\n", err)
	}
	return nil
}

func (cs *ConsulServer) DeregisterService(serviceName, ip string, port int) error {
	serviceID := fmt.Sprintf("%s-%s-%d", serviceName, ip, port)
	if err := cs.client.Agent().ServiceDeregister(serviceID); err != nil {
		return fmt.Errorf("failed to deregister service from consul: %s\n", err)
	}
	return nil
}
