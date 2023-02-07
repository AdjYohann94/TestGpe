package service

import (
	"fmt"
	"github.com/joho/godotenv"
	"strings"

	"github.com/spf13/viper"
)

const (
	KeyPostgresqlUser     = "postgres.user"
	KeyPostgresqlPassword = "postgres.password"
	KeyPostgresqlHost     = "postgres.host"
	KeyPostgresqlPort     = "postgres.port"
	KeyPostgresqlDBName   = "postgres.database"

	KeyJwtSecret = "jwt.secret"
	KeyLogLevel  = "log.level"
)

// LoadConfig load configuration file main.yaml and environment variables (from system and .env file) into the
// viper instance.
func LoadConfig(dir string) {
	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("fatal error during reading of config file .env: %w", err))
	}
	viper.SetConfigName("main")
	viper.AddConfigPath(dir + "/configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error during reading of config file: %w", err))
	}
}
