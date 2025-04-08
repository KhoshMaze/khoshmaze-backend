package config

import "time"

type Config struct {
	DB               DBConfig               `json:"database"`
	Cache            CacheConfig            `json:"cache"`
	Server           ServerConfig           `json:"server"`
	Jobs             JobsConfig             `json:"jobs"`
	SMS              SMSConfig              `json:"SMS"`
	Extra            ExtraConfig            `json:"extra"`
	AnomalyDetection AnomalyDetectionConfig `json:"anomalyDetection"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type CacheConfig struct {
	Redis            RedisConfig            `json:"redis"`
	MemoryCache      MemoryCacheConfig      `json:"memoryCache"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	ClientName string `json:"clientName"`
}

type MemoryCacheConfig struct {
	Size int `json:"size"`
	HealthCheckInterval time.Duration `json:"healthCheckInterval"`
}

type ServerConfig struct {
	Host              string `json:"host"`
	Port              uint   `json:"port"`
	AuthExpMinute     uint   `json:"authExpMin"`
	AuthRefreshMinute uint   `json:"authRefreshMin"`
	AuthSecret        string `json:"authSecret"`
	RefreshSecret     string `json:"refreshSecret"`
	AESSecret         string `json:"AESSecret"`
	SSLCertPath       string `json:"cert"`
	SSLKeyPath        string `json:"key"`
}

type AnomalyDetectionConfig struct {
	TTL         time.Duration `json:"ttl"`
	MaxSpeed    float64       `json:"maxSpeed"`
	MaxDistance float64       `json:"maxDistance"`
	DBPath      string        `json:"dbPath"`
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
	TemplateID string `json:"templateID"`
	OtpExpMin  int    `json:"otpExpMin"`
}

type ExtraConfig struct {
	FrontendDomain string `json:"frontendDomain"`
}
