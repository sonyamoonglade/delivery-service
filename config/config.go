package config

import (
	"github.com/spf13/viper"
)

func ReadConfig() (*viper.Viper, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return viper.GetViper(), nil
}
