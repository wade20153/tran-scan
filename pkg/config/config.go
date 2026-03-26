package config

// config 模块 — 解析 phase.yml
import (
	"os"

	"gopkg.in/yaml.v3"
)

type PhaseConfig struct {
	Crypto struct {
		Version string `yaml:"version"`
		System  struct {
			Key  string `yaml:"key"`
			Salt string `yaml:"salt"`
		} `yaml:"system"`
		Pbkdf2 struct {
			Iterations int `yaml:"iterations"`
			KeyLength  int `yaml:"keyLength"`
		} `yaml:"pbkdf2"`
		Aes struct {
			Mode string `yaml:"mode"`
		} `yaml:"aes"`
	} `yaml:"crypto"`
}

// LoadPhaseConfig 读取 phase.yml
func LoadPhaseConfig(path string) (*PhaseConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &PhaseConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
