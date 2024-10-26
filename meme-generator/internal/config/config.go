package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	envPrefix  = "MEME_"
	configName = "config.yaml"
	configType = "yaml"
)

type ConfigManager struct {
	Config *Config
	viper  *viper.Viper
}

func NewConfigManager() *ConfigManager {
	viper := viper.New()

	return &ConfigManager{
		Config: &Config{},
		viper:  viper,
	}
}

func (cm *ConfigManager) Load(configPath string) error {
	cm.viper.AddConfigPath(configPath)
	cm.viper.SetConfigName(configName)
	cm.viper.SetConfigType(configType)

	bindEnvs(cm.viper)

	setDefaults(cm.viper)

	if err := cm.readConfiguration(configName); err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	cm.viper.AutomaticEnv()

	if err := cm.viper.Unmarshal(&cm.Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}

func (cm *ConfigManager) readConfiguration(configFilePath string) error {
	err := cm.viper.ReadInConfig() // Find and read the config file

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, envPrefix) {
				return nil
			}
		}

		// if file does not exist, simply create one
		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			_, err := os.Create(configFilePath)
			if err != nil {
				return err
			}

			// write default config to the file
			err = cm.viper.WriteConfigAs(configFilePath)
			if err != nil {
				return err
			}
		}
	} else if err != nil {
		return err
	}

	return nil
}

func DebugMode() bool {
	if viper.IsSet("app.debug") {
		return viper.GetBool("app.debug")
	}

	return os.Getenv("DEBUG") == "true"
}
