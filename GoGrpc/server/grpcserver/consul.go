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

func (cs *ConsulServer) RegisterService(cfg *config.GrpcServerConfig, tags []string, check bool) error {
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s_%s_%d", cfg.Name, cfg.Addr, cfg.Port),
		Name:    cfg.ServerName,
		Tags:    tags,
		Address: cfg.Addr,
		Port:    cfg.Port,
	}

	if !check {
		fmt.Println("regist witchout check...")
		err := cs.client.Agent().ServiceRegister(registration)
		if err != nil {
			return fmt.Errorf("failed to register service to consul: %s\n", err)
		}
		return nil
	}

	fmt.Println("regist witch check...")
	registration.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Timeout:                        "3s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "600s",
	}
	if err := cs.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service to consul: %s\n", err)
	}
	return nil
}

func (cs *ConsulServer) DeregisterService(cfg *config.GrpcServerConfig) error {
	if err := cs.client.Agent().ServiceDeregister(fmt.Sprintf("%s_%s_%d", cfg.Name, cfg.Addr, cfg.Port)); err != nil {
		return fmt.Errorf("failed to deregister service from consul: %s\n", err)
	}
	return nil
}
