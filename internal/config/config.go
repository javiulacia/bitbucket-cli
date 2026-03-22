package config

type AuthMethod string

const (
	AuthMethodAPIToken    AuthMethod = "api_token"
	AuthMethodAccessToken AuthMethod = "access_token"
)

type Config struct {
	BaseURL     string     `json:"base_url"`
	Workspace   string     `json:"workspace,omitempty"`
	AuthMethod  AuthMethod `json:"auth_method"`
	Email       string     `json:"email,omitempty"`
	APIToken    string     `json:"api_token,omitempty"`
	AccessToken string     `json:"access_token,omitempty"`
	GitProtocol string     `json:"git_protocol,omitempty"`
}
