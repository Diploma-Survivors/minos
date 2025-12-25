package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Gemini   GeminiConfig
}

type ServerConfig struct {
	Port      string
	ApiPrefix string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type GeminiConfig struct {
	ApiKey string
	Model  string
}

func NewConfig() (*Config, error) {
	// Configure Viper to read .env file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Enable automatic environment variable loading
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Warn().Err(err).Msg("Error reading config file")
	}

	var config Config
	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Server.ApiPrefix = viper.GetString("API_PREFIX")
	if config.Server.ApiPrefix == "" {
		config.Server.ApiPrefix = "/api/v1"
	}
	config.Database.Host = viper.GetString("DATABASE_HOST")
	config.Database.Port = viper.GetString("DATABASE_PORT")
	config.Database.User = viper.GetString("DATABASE_USER")
	config.Database.Password = viper.GetString("DATABASE_PASSWORD")
	config.Database.Name = viper.GetString("DATABASE_NAME")
	config.Redis.Host = viper.GetString("REDIS_HOST")
	config.Redis.Port = viper.GetString("REDIS_PORT")
	config.Redis.Password = viper.GetString("REDIS_PASSWORD")
	config.Redis.DB = viper.GetInt("REDIS_DB")
	config.Gemini.ApiKey = viper.GetString("GEMINI_API_KEY")
	config.Gemini.Model = viper.GetString("GEMINI_MODEL")

	log.Info().Interface("config", config).Msg("Config loaded")
	return &config, nil
}
