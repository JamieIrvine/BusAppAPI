package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Ingestion Ingestion `mapstructure:"auth"`
}

type Ingestion struct {
	GTFS []GTFSIngestion `mapstructure:"gtfs"`
}

type GTFSIngestion struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

func InitConfig() (Config, error) {
	configFilePath := "./config/cfg.yaml"
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv() // Enable automatic environment variable binding
	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	// Override config values with environment variables
	_ = viper.BindEnv("database.username", "POSTGRES_USER")
	_ = viper.BindEnv("database.password", "POSTGRES_PASSWORD")
	_ = viper.BindEnv("database.host", "POSTGRES_HOST")
	_ = viper.BindEnv("database.databaseName", "POSTGRES_DB_NAME")

	var appConfig Config
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return Config{}, err
	}

	log.Printf("Running with config: %s\n", configFilePath)

	return appConfig, nil
}
