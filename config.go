package main

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/Yeremi528/itudy-back/kit/secretmanager"
	"github.com/spf13/viper"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Web         WebConfig    `yaml:"API"`
	Debug       DebugConfig  `yaml:"DEBUG"`
	MercadoPago MPConfig     `yaml:"MERCADO_PAGO"`
	Resend      ResendConfig `yaml:"RESEND"`
	Mongo       MongoConfig  `yaml:"MONGO"`
}

type WebConfig struct {
	Host            string        `validate:"required" yaml:"HOST"`
	ReadTimeout     time.Duration `validate:"required" yaml:"READ_TIMEOUT"`
	IdleTimeout     time.Duration `validate:"required" yaml:"IDLE_TIMEOUT"`
	WriteTimeout    time.Duration `validate:"required" yaml:"WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `validate:"required" yaml:"SHUTDOWN_TIMEOUT"`
}
type DebugConfig struct {
	Host            string        `validate:"required" yaml:"HOST"`
	ReadTimeout     time.Duration `validate:"required" yaml:"READ_TIMEOUT"`
	IdleTimeout     time.Duration `validate:"required" yaml:"IDLE_TIMEOUT"`
	WriteTimeout    time.Duration `validate:"required" yaml:"WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `validate:"required" yaml:"SHUTDOWN_TIMEOUT"`
}

type ResendConfig struct {
	APIKEY string `yaml:"APIKEY"`
}

type MPConfig struct {
	AccessToken string `yaml:"ACCESS_TOKEN"`
}

type MongoConfig struct {
	Conexion string `yaml:"CONEXION"`
}

func loadConfig(ctx context.Context, sm *secretmanager.Client) (Config, error) {
	var config Config
	var configData []byte
	var err error

	env := os.Getenv("env")
	if env == "local" {
		// Cargar configuración desde archivo local
		configData, err = os.ReadFile("config.yaml")

		if err != nil {
			return Config{}, err
		}
	} else {
		// Cargar desde el secret manager
		secret, err := sm.GetSecret(ctx, "ITUDY")
		if err != nil {
			return Config{}, err
		}
		configData = secret
	}

	// Parsear YAML (tanto para local como para secret manager)
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

// LoadConfig reads configuration from file or environment variables.
func readConfig() (Config, error) {
	viper.AddConfigPath("./app/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, err
	}

	// if err := modelvalidator.Check(&cfg, false); err != nil {
	// 	return Config{}, err
	// }

	return cfg, nil
}
