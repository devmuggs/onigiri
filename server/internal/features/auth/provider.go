package auth

import "fmt"

func GetProviderConfig(providerName string, cfg []ProviderConfig) (*ProviderConfig, error) {
	for _, p := range cfg {
		if p.Name == providerName {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("provider %s not found", providerName)
}
