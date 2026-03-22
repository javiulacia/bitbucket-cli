package bitbucket

import (
	"encoding/base64"
	"fmt"

	"bitbucket-cli/internal/config"
)

func authHeader(cfg *config.Config) (string, error) {
	switch cfg.AuthMethod {
	case config.AuthMethodAPIToken:
		if cfg.Email == "" || cfg.APIToken == "" {
			return "", fmt.Errorf("api token auth requires both email and api token")
		}
		token := base64.StdEncoding.EncodeToString([]byte(cfg.Email + ":" + cfg.APIToken))
		return "Basic " + token, nil
	case config.AuthMethodAccessToken:
		if cfg.AccessToken == "" {
			return "", fmt.Errorf("access token auth requires an access token")
		}
		return "Bearer " + cfg.AccessToken, nil
	default:
		return "", fmt.Errorf("unsupported auth method %q", cfg.AuthMethod)
	}
}
