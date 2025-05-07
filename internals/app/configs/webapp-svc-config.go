package appconfigs

import (
	filepath "path/filepath"
)

type WebAppConfig struct {
	Path         string
	AuthnSecrets struct {
		Path            string `yaml:"path"`
		JWTAuthnSecrets struct {
			Path string `yaml:"path"`
		} `yaml:"jwtAuthnSecrets"`
	} `yaml:"authnSecrets"`
	NetworkConfigFile string `yaml:"networkConfigFile"`
	URIs              struct {
		Login    string `yaml:"login" default:"auth/login"`
		Logout   string `yaml:"logout" default:"auth/logout"`
		Callback string `yaml:"callback" default:"auth/callback"`
	} `yaml:"uris"`
}

func (sc *WebAppConfig) GetAuthnSecretPath() string {
	return filepath.Join(sc.Path, sc.AuthnSecrets.Path)
}

func (sc *WebAppConfig) GetJwtAuthSecretPath() string {
	return filepath.Join(sc.Path, sc.AuthnSecrets.JWTAuthnSecrets.Path)
}
