package auth

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*AuthConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	var config AuthConfig
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
