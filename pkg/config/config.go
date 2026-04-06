package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Name     string
	Username string
	Password string
}

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  int
	RefreshTokenExpiry int
}

type AppConfig struct {
	Version    string
	APIVersion string
	Locale     string
	SiteOwner  string
	AppEnv     string
	AppURL     string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("DB_CONNECTION", "mysql")
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("JWT_ACCESS_EXPIRY_HOURS", 24)
	viper.SetDefault("JWT_REFRESH_EXPIRY_DAYS", 14)
	viper.SetDefault("APP_VERSION", "6.5.9")
	viper.SetDefault("API_VERSION", "2.1.0")
	viper.SetDefault("DEFAULT_LOCALE", "en_US")
	viper.SetDefault("APP_ENV", "production")
	viper.SetDefault("APP_URL", "http://localhost:8080")

	_ = viper.ReadInConfig()

	cfg := &Config{}
	cfg.Server.Port = viper.GetString("PORT")
	cfg.Server.Mode = viper.GetString("GIN_MODE")
	cfg.Database.Driver = viper.GetString("DB_CONNECTION")
	cfg.Database.Host = viper.GetString("DB_HOST")
	cfg.Database.Port = viper.GetString("DB_PORT")
	cfg.Database.Name = viper.GetString("DB_DATABASE")
	cfg.Database.Username = viper.GetString("DB_USERNAME")
	cfg.Database.Password = viper.GetString("DB_PASSWORD")
	cfg.JWT.Secret = viper.GetString("JWT_SECRET")
	cfg.JWT.AccessTokenExpiry = viper.GetInt("JWT_ACCESS_EXPIRY_HOURS")
	cfg.JWT.RefreshTokenExpiry = viper.GetInt("JWT_REFRESH_EXPIRY_DAYS")
	cfg.App.Version = viper.GetString("APP_VERSION")
	cfg.App.APIVersion = viper.GetString("API_VERSION")
	cfg.App.Locale = viper.GetString("DEFAULT_LOCALE")
	cfg.App.SiteOwner = viper.GetString("SITE_OWNER")
	cfg.App.AppEnv = viper.GetString("APP_ENV")
	cfg.App.AppURL = viper.GetString("APP_URL")

	return cfg, nil
}
