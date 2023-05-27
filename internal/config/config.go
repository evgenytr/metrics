package config

type AgentConfig struct {
	Host           string  `env:"ADDRESS"`
	ReportInterval float64 `env:"REPORT_INTERVAL"`
	PollInterval   float64 `env:"POLL_INTERVAL"`
}

type ServerConfig struct {
	Host string `env:"ADDRESS"`
}
