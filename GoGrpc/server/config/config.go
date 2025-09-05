package config

import (
	"fmt"
	"server/utils"
)

type Config struct {
	GrpcServers []GrpcServerConfig
	Consul      ConsulConfig
}

type GrpcServerConfig struct {
	ServerName string
	Name       string
	Addr       string
	Port       int
}

type ConsulConfig struct {
	Scheme string // http or https
	Addr   string // 127.0.0.1:8500
}

func DefaultConfig() *Config {
	outboundIP, err := utils.GetOutboundIP()
	if err != nil {
		fmt.Printf("Failed to get outbound IP: %s\n", err)
		return nil
	}
	fmt.Println("Outbound IP:", outboundIP)
	return &Config{
		GrpcServers: []GrpcServerConfig{
			{
				ServerName: "grpcServer",
				Name:       "grpcServer1",
				Addr:       outboundIP,
				Port:       50051,
			},
			{
				ServerName: "grpcServer",
				Name:       "grpcServer2",
				Addr:       outboundIP,
				Port:       50052,
			},
		},
		Consul: ConsulConfig{
			Scheme: "http",
			Addr:   fmt.Sprintf("%s:%d", outboundIP, 8500),
		},
	}
}
