package config

type Config struct {
	DB     DBConfig     `json:"database"`
	Redis  RedisConfig  `json:"redis"`
	Server ServerConfig `json:"server"`
	Jobs   JobsConfig   `json:"jobs"`
	SMS    SMSConfig    `json:"SMS"`
	Extra  ExtraConfig  `json:"extra"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
}

type ServerConfig struct {
	Host              string `json:"host"`
	Port              uint   `json:"port"`
	AuthExpMinute     uint   `json:"authExpMin"`
	AuthRefreshMinute uint   `json:"authRefreshMin"`
	AuthSecret        string `json:"authSecret"`
	RefreshSecret     string `json:"refreshSecret"`
	SSLCertPath       string `json:"cert"`
	SSLKeyPath        string `json:"key"`
}

type JobsConfig struct {
	TokenCheckerIntervalMinute int `json:"tokenCheckerInterval"`
	OutboxPollerIntervalSecond int `json:"outboxPollerInterval"`
}

type SMSConfig struct {
	ApiKey       string                `json:"apiKey"`
	ApiBaseURL   string                `json:"apiBaseURL"`
	Sender       string                `json:"sender"`
	Verification SMSVerificationConfig `json:"verification"`
}

type SMSVerificationConfig struct {
	URL        string `json:"url"`
	TemplateID string    `json:"templateID"`
	OtpExpMin  int    `json:"otpExpMin"`
}

type ExtraConfig struct {
	FrontendDomain string `json:"frontendDomain"`
}
