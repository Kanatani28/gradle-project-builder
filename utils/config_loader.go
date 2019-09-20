package utils

import (
	"log"

	"github.com/spf13/viper"
)

// config.yml内のprojectsの１要素を表す構造体
type Project struct {
	Name           string
	RepositoryPath string
	OutputName     string
	BuildType      string
	Ignore         bool
}

// config.ymlを表す構造体
type Config struct {
	Projects    []Project
	BuildTarget string
	OutputDir   string
}

const configName = "config"
const configType = "yml"

// main.goがある階層から
const configDir = "."

func LoadConfig() Config {

	var C Config

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDir)
	// viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%v", err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatalf("%v", err)
	}
	// log.Printf("%+v", C)

	return C
}

// application.ymlを表す構造体
type ApplicationYml struct {
	Spring struct {
		Profiles struct {
			Active string
		}
	}
}

const (
	appConfigName = "application"
	appConfigExt  = "yml"
)

func ReplaceAppConfig(confPath string, env string) {

	v := viper.New()
	log.Println("Start replace app config.")

	v.SetConfigName(appConfigName)
	v.SetConfigType(appConfigExt)
	v.AddConfigPath(confPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("%v", err)
	}
	// fmt.Printf("%+v\n", viper.AllSettings())
	v.Set("spring.profiles.active", env)
	v.WriteConfig()

	log.Println("Finishsed replace app config . " + env)

}
