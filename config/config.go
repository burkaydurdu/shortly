package config

import "fmt"

/*
 *
 * Set Environment Variables,
 * We may use Viper[github.com/spf13/viper] and pp[github.com/k0kubun/pp] 🚀
 *
 */
const (
	AppName = "Shortly"
	Port    = 6161
	IsDebug = true
)

type ServerConfig struct {
	Port int
}

type Config struct {
	AppName string
	IsDebug bool
	Server  ServerConfig
}

func New() (*Config, error) {
	config := &Config{}

	config.AppName = AppName
	config.IsDebug = IsDebug
	config.Server = ServerConfig{
		Port: Port,
	}

	return config, nil
}

func (c Config) Print() {
	fmt.Printf("%+v\n", c) // [TODO] should improve in here
}