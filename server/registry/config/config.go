package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MongodbURI      string `mapstructure:"mongodb_uri"`
	UserName        string `mapstructure:"user_name"`
	UserPassword    string `mapstructure:"user_password"`
	AdminName       string `mapstructure:"admin_name"`
	AdminPassword   string `mapstructure:"admin_name"`
	UserServerPort  string `mapstructure:"user_port"`
	AdminServerPort string `mapstructure:"admin_port"`
}

var AppConfig Config

func LoadAppConfig() {
	log.Println("Loading Server Configurations...")
	viper.AddConfigPath("/home/jean/")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		log.Fatal(err)
	}
}
