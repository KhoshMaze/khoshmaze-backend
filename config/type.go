package config

type Config struct {
	DB     DBConfig     `json:"database"`
	Redis  RedisConfig  `json:"redis"`
	Server ServerConfig `json:"server"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

type ServerConfig struct {
	Host             string `json:"host"`
	Port             uint   `json:"port"`
	AuthExpMinute    uint   `json:"authExpMin"`
	AuthRefreshMinut uint   `json:"authRefreshMin"`
	RefreshSecret    string `json:"refreshSecret"`
	Salt             string `json:"salt"`
}
