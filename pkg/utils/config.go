package utils

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Database DatabaseSettings `yaml:"database"`
	Server   ServerSettings   `yaml:"server"`
	Redis    RedisSettings    `yaml:"redis"`
}

type DatabaseSettings struct {
	Url    string `yaml:"url"`
	DbName string `yaml:"dbName"`
}

type RedisSettings struct {
	Address string `yaml:"address"`
}

type ServerSettings struct {
	Port string `yaml:"port"`
}

func ReadConfig() (Config, error) {
	
	databaseConfig, err := readDatabaseConfig()

	if err != nil {
		return databaseConfig, err
	}

	appConfig, err := readAppConfig()

	if err != nil {
		return appConfig, err
	}

	config := Config{Database: databaseConfig.Database, Redis: databaseConfig.Redis, Server: appConfig.Server}

	return config, nil

}

func readDatabaseConfig() (Config, error) {
	var config Config

	config.Database.Url = os.Getenv("MONGO_URL")
	config.Database.DbName = "ideanest"
	config.Redis.Address = os.Getenv("REDIS_ADDRESS")
	
	if config.Database.Url == ""{

		viper.SetConfigName("database-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		if err := viper.ReadInConfig(); err != nil {
			return config, err
		}

		if err := viper.Unmarshal(&config); err != nil {
			return config, err
		}
		
	}
	
	return config, nil
}

func readAppConfig() (Config, error) {
	var config Config

	config.Server.Port = os.Getenv("PORT")
	
	if config.Server.Port == ""{

		viper.SetConfigName("app-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./config")

		if err := viper.ReadInConfig(); err != nil {
			return config, err
		}

		if err := viper.Unmarshal(&config); err != nil {
			return config, err
		}
		
		
	}
	
	
	return config, nil
}
