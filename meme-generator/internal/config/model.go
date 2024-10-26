package config

type Config struct {
	App      AppConfiguration      `mapstructure:"app"`
	Web      WebConfiguration      `mapstructure:"web"`
	Database DatabaseConfiguration `mapstructure:"database"`
}

type AppConfiguration struct {
	Debug   bool   `mapstructure:"debug"`
	DataDir string `mapstructure:"data_dir"`
}

type WebConfiguration struct {
	Port      int    `mapstructure:"port"`
	SecretKey string `mapstructure:"secret_key"`
	Static    string `mapstructure:"static"`
}

type DatabaseConfiguration struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	LogMode  bool   `mapstructure:"log_mode"`
	SSLMode  string `mapstructure:"ssl_mode"`
}
