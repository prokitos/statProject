package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// структура конфигов, а также их считывание из файла, и загрузка в эти стуктуры.

type MainConfig struct {
	Server     ServerConfig   `yaml:"server"`
	PostgresDB PostgresConfig `yaml:"postgres"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}
type PostgresConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Host   string `yaml:"host"`
	Name   string `yaml:"name"`
	Reload bool   `yaml:"reload"`
	Port   string `yaml:"port"`
}

func (cfg *MainConfig) ConfigMustLoad(name string) {

	path := "./config/" + name + ".yaml"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		pathBackup := "../config/" + name + ".yaml"
		if _, err := os.Stat(pathBackup); os.IsNotExist(err) {
			panic("nothing from this path")
		}
		path = pathBackup
	}

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read")
	}

}
