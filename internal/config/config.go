package config

import (
	"fmt"
	"os"
)

type Config struct {
	TopEndpoint string
}

var (
	Cfg Config
)

const (
	DefaultTopEndpoint = "localhost:30000"
)

func Print() {
	fmt.Printf("Top endpoint: %s\n", Cfg.TopEndpoint)
}

func init() {
	Cfg = Config{
		TopEndpoint: getEnv("TOP_ENDPOINT", DefaultTopEndpoint),
	}
}

func getEnv(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return value
}
