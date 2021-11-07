package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"log"
	"rest_api/internal/app/apiserver"
)

var(
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}