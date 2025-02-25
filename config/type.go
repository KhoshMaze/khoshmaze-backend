package config

type Config struct {
	DB     DBConfig     `json:"database"`
	Redis  RedisConfig  `json:"redis"`
	Server ServerConfig `json:"server"`
	Jobs   JobsConfig   `json:"jobs"`
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
}

type JobsConfig struct {
	TokenCheckerInterval int `json:"tokenCheckerInterval"`
}
