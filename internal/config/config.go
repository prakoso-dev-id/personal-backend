package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port         string
	Mode         string
	BaseURL      string `mapstructure:"base_url"`
	CleanStorage bool
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int // hours
}

func LoadConfig() (*Config, error) {
	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.base_url", "http://localhost:8080")
	viper.SetDefault("jwt.expiration", 24)
	viper.SetDefault("database.sslmode", "disable")

	// Read config.yaml first
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config.yaml not found, using defaults/environment variables")
		} else {
			log.Printf("Warning: error reading config.yaml: %v", err)
		}
	}

	// Try to read .env file (optional, won't crash if missing)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	if err := viper.MergeInConfig(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// Always read OS environment variables (Docker injects these)
	viper.AutomaticEnv()

	// Map ENV vars to config struct
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("server.base_url", "SERVER_BASE_URL")

    // Manual mapping for flat .env file to nested struct
    if viper.IsSet("DB_HOST") { viper.Set("database.host", viper.GetString("DB_HOST")) }
    if viper.IsSet("DB_PORT") { viper.Set("database.port", viper.GetString("DB_PORT")) }
    if viper.IsSet("DB_USER") { viper.Set("database.user", viper.GetString("DB_USER")) }
    if viper.IsSet("DB_PASSWORD") { viper.Set("database.password", viper.GetString("DB_PASSWORD")) }
    if viper.IsSet("DB_NAME") { viper.Set("database.dbname", viper.GetString("DB_NAME")) }
    // SSLMode might not be in .env, default was set above
    if viper.IsSet("DB_SSLMODE") { viper.Set("database.sslmode", viper.GetString("DB_SSLMODE")) }

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
