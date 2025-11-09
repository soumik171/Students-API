package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Addr string `yaml:"address" env-required:"true"`
}

// struct tags:

// env-default:"production"

type Config struct {
	Env        string `yaml:"env" env:"ENV" env-required:"true" `
	Storage    string `yaml:"storage_path" env-required:"true"`
	HttpServer `yaml:"http_server"`
}

// Inside MustLoad function->

// 1. search for the config path
// 2.if not then check with the flag
// 3.Also if not, then "Not Found"
// 4.After getting path, if there's an error, then check it's actually error or missing files or something else(at once with os.Stat() & os.IsNotExist())
// 5. Read the file using cleanenv, if not then print file cannot read

// log.Fatal()-> logs error & stops the program
// log.Fatalf()-> logs error & stops the program & also format that

func MustLoad() *Config {

	var configPath string

	// 1.
	configPath = os.Getenv("CONFIG_PATH")

	// 2.
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file") // Declares the flag
		flag.Parse()                                                         // Fill the flag with actual value

		configPath = *flags

		// 3.
		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// 4.

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	// 5.

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg) // ReadConfig(config location, memory location of where to store that)

	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}

	return &cfg

}
