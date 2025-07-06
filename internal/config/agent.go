package config

type AgentConfig struct {
	Endpoint string
	Secret   string
}

func NewAgentConfig(config *Config) *AgentConfig {
	return &AgentConfig{
		Endpoint: config.AgentEndpoint,
		Secret:   config.AgentSecret,
	}
}
