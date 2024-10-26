package config

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.debug", false)
	v.SetDefault("app.data_dir", "./data")

	v.SetDefault("web.port", 8080)
	v.SetDefault("web.secret_key", "secret")
	v.SetDefault("web.static", "./static")

	v.SetDefault("database.name", "memes")
	v.SetDefault("database.username", "memes")
	v.SetDefault("database.password", "memes")
	v.SetDefault("database.host", "db")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.log_mode", false)
	v.SetDefault("database.ssl_mode", "disable")
}
