package config

import "github.com/spf13/viper"

func bindEnvs(v *viper.Viper) {
	v.BindEnv("app.debug", "MEME_APP_DEBUG")
	v.BindEnv("app.data_dir", "MEME_APP_DATA_DIR")

	v.BindEnv("web.port", "MEME_WEB_PORT")
	v.BindEnv("web.secret_key", "MEME_WEB_SECRET_KEY")
	v.BindEnv("web.static", "MEME_WEB_STATIC")

	v.BindEnv("database.name", "MEME_DATABASE_NAME")
	v.BindEnv("database.username", "MEME_DATABASE_USERNAME")
	v.BindEnv("database.password", "MEME_DATABASE_PASSWORD")
	v.BindEnv("database.host", "MEME_DATABASE_HOST")
	v.BindEnv("database.port", "MEME_DATABASE_PORT")
	v.BindEnv("database.log_mode", "MEME_DATABASE_LOG_MODE")
	v.BindEnv("database.ssl_mode", "MEME_DATABASE_SSL_MODE")
}
