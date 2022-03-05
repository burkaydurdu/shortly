package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

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
	AppName               string
	IsDebug               bool
	LengthOfCode          int
	DurationOfWriteToDisk time.Duration
	Server                ServerConfig
	MemoryPath            string
	MemoryFileName        string
}

func New() (*Config, error) {
	config := &Config{}

	config.AppName = shortlyViberString("NAME", AppName)
	config.IsDebug = IsDebug
	config.LengthOfCode = 6
	config.DurationOfWriteToDisk = time.Second * 2
	config.MemoryPath = ".mem"
	config.MemoryFileName = "shortly"
	config.Server = ServerConfig{
		Port: shortlyViberInt("PORT", Port),
	}

	return config, nil
}

func (c *Config) Print() {
	fmt.Printf("%+v\n", c) // [TODO] should improve in here
}

func shortlyViberString(envKey, defaultKey string) string {
	envValue := os.Getenv(envKey)

	if len(envValue) == 0 {
		return defaultKey
	}

	return envValue
}

func shortlyViberInt(envKey string, defaultValue int) int {
	envValue := os.Getenv(envKey)

	envIntegerValue, err := strconv.ParseInt(envValue, 10, 64)

	if err != nil {
		return defaultValue
	}

	return int(envIntegerValue)
}
