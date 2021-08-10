package util

import (
	"time"

	"github.com/spf13/viper"
)

var (
	TZLocation   *time.Location
	DATABASE_URL string
)

func init() {
	ReadEnv()
	timezone := viper.GetString("TIMEZONE")
	TZLocation, _ = time.LoadLocation(timezone)
	DATABASE_URL = viper.GetString("DATABASE_URL")
}

func ReadEnv() error {
	viper.SetConfigFile(".env")
	return viper.ReadInConfig()
}
