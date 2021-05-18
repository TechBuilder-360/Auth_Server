package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)

//Data : config data
type Data struct {
	ServerHost string `yaml:"ServerHost"`
	ServerPort string `yaml:"ServerPort"`
	SecretKey  string `yaml:"SecretKey"`
	PGUser     string `yaml:"PGUser"`
	PGDatabase string `yaml:"PGDatabase"`
	PGPassword string `yaml:"PGPassword"`
	PGHost     string `yaml:"PGHost"`
	PGPort     string `yaml:"PGPort"`
}

//Init : initialize data
func (c *Data) Init(configDir string) {

	dir, dirErr := os.Getwd()
	if dirErr != nil {
		log.Printf("Cannot set default input/output directory to the current working directory >> %s", dirErr)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(dir)
	viper.AddConfigPath(configDir)
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("\n fatal error: could not read from config file >>%s ", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		err := viper.ReadInConfig()
		if err != nil {
			log.Printf("\n fatal error: could not read from config file >>%s ", err)
		}
		viper.Unmarshal(c)
	})

	viper.Unmarshal(c)
	log.Println("App configuration loaded successfully!")
}
