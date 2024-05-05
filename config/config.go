package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	RefreshTokenExpiration = time.Hour * 24 * 7
	AccessTokenExpiration  = time.Minute * 10
)

type jwtConfig struct {
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenSecret  string `mapstructure:"ACCES_TOKEN_SECRET"`
}

// config defines the structure for the application configuration.
// It includes the server port and the data source name (DSN) for database connection.
type config struct {
	PORT int       `mapstructure:"port"` // PORT defines the port on which the server should run.
	DSN  string    `mapstructure:"dsn"`  // DSN is the Data Source Name for the database connection.
	JWT  jwtConfig `mapstructure:"JWT"`
}

// cfg holds the application configuration loaded from the config file.
var cfg config

// loadConfig reads the configuration from a YAML file and unmarshals it into the cfg variable.
// It looks for the configuration file named "config" with a ".yaml" extension in the current directory
// and the "/app/" directory. If reading or unmarshaling fails, the function panics.
func loadConfig() {
	viper.SetConfigName("config") // Specifies the name of the config file to look for.
	viper.SetConfigType("yaml")   // Sets the format of the config file.
	viper.AddConfigPath(".")      // Adds the current directory as a path to look for the config file.
	viper.AddConfigPath("/app/")  // Adds the "/app/" directory as a path to look for the config file.

	// Reads the config file and checks for errors.
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	// Unmarshals the read config into the cfg variable and checks for errors.
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
}

// Config returns the application configuration.
// It ensures that the configuration is loaded before returning it.
// If the configuration has not been loaded yet (indicated by an empty DSN or PORT),
// it calls loadConfig to load the configuration.
func Config() config {
	if cfg.DSN == "" || cfg.PORT == 0 {
		loadConfig() // Loads the configuration if it's not already loaded.
	}
	return cfg // Returns the loaded configuration.
}
