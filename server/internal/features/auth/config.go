package auth

type ProviderConfig struct {
	Name          string            `yaml:"name"`
	Type          string            `yaml:"type,omitempty"` // for "local", etc
	ClientID      string            `yaml:"client_id"`
	ClientSecret  string            `yaml:"client_secret"`
	AuthURL       string            `yaml:"auth_url"`
	TokenURL      string            `yaml:"token_url"`
	UserInfoURL   string            `yaml:"user_info_url"`
	Scopes        []string          `yaml:"scopes"`
	FieldMappings map[string]string `yaml:"field_mappings"`
	Transforms    []TransformConfig `yaml:"transforms"`
}

type TransformConfig struct {
	Name   string            `yaml:"name"`
	Params map[string]string `yaml:"params"`
}

type AuthConfig struct {
	Providers []ProviderConfig `yaml:"providers"`
}
