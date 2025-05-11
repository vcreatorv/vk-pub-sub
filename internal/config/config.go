package config

import "fmt"

type AppConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (app *AppConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", app.Host, app.Port)
}

type SubPubConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (sub *SubPubConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", sub.Host, sub.Port)
}
