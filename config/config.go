package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Interval int `hcl:"interval"`
	DDNS     struct {
		Service []Service `hcl:"service,block"`
	} `hcl:"ddns,block"`
}

type Service struct {
	Provider   string `hcl:"provider,label"`
	APIKey     string `hcl:"api_key"`
	ZoneID     string `hcl:"zone_id"`
	RecordName string `hcl:"record_name"`
	TTL        int    `hcl:"ttl"`
}

func LoadConfig() (*Config, error) {
	cfg := new(Config)
	if err := hclsimple.DecodeFile("./config.hcl", nil, cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
