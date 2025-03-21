package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Url1          string
	Url2          string
	TimeoutMillis int64
	Cep           string
}

func init() {
	log.Println("Loading config file...")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println("Config file loaded successfully!")
}

func GetConfig() *config {
	c := &config{}
	err := viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}
	return c
}

func (c *config) GetUrl1(cep string) string {
	r := fmt.Sprintf(c.Url1, cep)
	return r
}

func (c *config) GetUrl2(cep string) string {
	return fmt.Sprintf(c.Url2, cep)
}
