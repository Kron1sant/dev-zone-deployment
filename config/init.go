package config

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

func init() {
	// Find and set config file
	setConfigFile()

	// Read config from file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error read config file: %v\n", err)
	}

	// Set a prefix to search for environment variables
	viper.SetEnvPrefix("DEVZONE")
	viper.AutomaticEnv()

	// Init all configs
	appConfInit()
	dbConfInit()
	ovpnConfInit()
	yacloudConfInit()
}

// ToDo
// func Check() {
// 	if c.appSecret == "" {
// 		log.Printf("App secret is not set. Use the config file or Env")
// 	}
// 	if c.dbPassword == "" {
// 		log.Printf("DataBase password is not set. Use the config file or Env")
// 	}
// 	if c.adminPassword == "" {
// 		log.Printf("Initial Admin password is not set. Use the config file or Env")
// 	}
// 	if c.yaCloudToken == "" {
// 		log.Printf("Yandex cloud token is not set. Use the config file or Env")
// 	}
// }

func setConfigFile() {
	defaultConfig := "config.default.yaml" // expected in the app directory
	configFileFromCLI := getConfigFileFromCLI()
	if _, err := os.Stat(configFileFromCLI); err != nil {
		viper.SetConfigFile(defaultConfig)
	} else {
		viper.SetConfigFile(configFileFromCLI)
	}
}

func getConfigFileFromCLI() string {
	// Disable error checking, since we parse only one key,
	// otherwise, if several keys were passed at startup,
	// we would get an error
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()
	return *configPath
}
