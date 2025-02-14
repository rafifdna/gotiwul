package config

import (
    "os"
    "time"
    
    "gopkg.in/yaml.v2"
)

type Config struct {
    Server struct {
        Port     int      `yaml:"port"`
        CertFile string   `yaml:"cert_file"`
        KeyFile  string   `yaml:"key_file"`
        Domain   string   `yaml:"domain"`
    } `yaml:"server"`
    
    Health struct {
        Interval time.Duration `yaml:"interval"`
        Timeout  time.Duration `yaml:"timeout"`
        Path     string        `yaml:"path"`
    } `yaml:"health"`
    
    Backends []string `yaml:"backends"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}