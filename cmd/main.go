package main

import (
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ositlar/effective/internal/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-flag", "config/config.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	os.Setenv("BIND_ADDR", config.BindAddr)
	os.Setenv("LOG_LEVEL", config.LogLevel)
	os.Setenv("DB_URL", config.DatabaseURL)
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	if err := apiserver.StartServer(config); err != nil {
		log.Fatal(err)
	}
}
