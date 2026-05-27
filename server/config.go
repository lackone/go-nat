package main

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenerFile string `yaml:"listener_file"`
}

func ParseConfig(file string) (*Config, error) {
	all, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(all, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type ListenerConfig struct {
	ClientId         string `json:"client_id"`         // 客户端ID
	PublicIp         string `json:"public_ip"`         // 公共IP
	PublicPort       int    `json:"public_port"`       // 公共端口
	PublicProtocol   string `json:"public_protocol"`   // 公共协议
	InternalIp       string `json:"internal_ip"`       // 内网IP
	InternalPort     int    `json:"internal_port"`     // 内网端口
	InternalProtocol string `json:"internal_protocol"` // 内网协议
}

func ParseListenerConfig(file string) ([]*ListenerConfig, error) {
	all, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	cfg := make([]*ListenerConfig, 0)
	err = json.Unmarshal(all, &cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
