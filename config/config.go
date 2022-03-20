package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var AppConfig *viper.Viper

type secConfig struct {
	appSecret     string
	dbPassword    string
	adminPassword string
	yaCloudToken  string
}

var SecConfig secConfig

func init() {
	flag.CommandLine.Init("CLI", flag.ContinueOnError)
	pConfigPath := flag.String("config", "", "Path to config file")
	//pTestMode := flag.String("test.run", "", "test mode")
	flag.Parse()

	// if *pTestMode != "" {
	// 	log.Println("Running in test mode")
	// }

	viper.SetEnvPrefix("DEVZONE")
	viper.AutomaticEnv()

	if *pConfigPath == "" {
		if viper.IsSet("CONFIG") {
			*pConfigPath = viper.GetString("CONFIG")
		} else {
			*pConfigPath = "config.yaml"
		}
	}

	viper.SetConfigFile(configFile(*pConfigPath))

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error read config file: %v\n", err)
	}
	viper.SetDefault("service.port", "8088")

	if viper.IsSet("APP_SECRET") {
		SecConfig.appSecret = viper.GetString("APP_SECRET")
	} else {
		SecConfig.appSecret = viper.GetString("service.secret")
	}
	if viper.IsSet("DB_PASSWORD") {
		SecConfig.dbPassword = viper.GetString("DB_PASSWORD")
	} else {
		SecConfig.dbPassword = viper.GetString("db.pass")
	}
	if viper.IsSet("ADMIN_PASSWORD") {
		SecConfig.adminPassword = viper.GetString("ADMIN_PASSWORD")
	} else {
		SecConfig.adminPassword = viper.GetString("service.default_admin.password")
	}
	if viper.IsSet("YA_TOKEN") {
		SecConfig.yaCloudToken = viper.GetString("YA_TOKEN")
	} else {
		SecConfig.yaCloudToken = viper.GetString("yacloud.token")
	}
	SecConfig.Check()

	AppConfig = viper.GetViper()

	if viper.IsSet("PORT") {
		AppConfig.Set("service.port", viper.GetString("PORT"))
	}
}

func (c *secConfig) GetAppSecret() string {
	return c.appSecret
}

func (c *secConfig) GetDBPassword() string {
	return c.dbPassword
}

func (c *secConfig) GetAdminPassword() string {
	return c.adminPassword
}

func (c *secConfig) GetYaCloudToken() string {
	return c.yaCloudToken
}

func (c *secConfig) Check() {
	if c.appSecret == "" {
		log.Printf("App secret is not set. Use the config file or Env")
	}
	if c.dbPassword == "" {
		log.Printf("DataBase password is not set. Use the config file or Env")
	}
	if c.adminPassword == "" {
		log.Printf("Initial Admin password is not set. Use the config file or Env")
	}
	if c.yaCloudToken == "" {
		log.Printf("Yandex cloud token is not set. Use the config file or Env")
	}
}

func configFile(pathToConfig string) string {
	stat, err := os.Stat(pathToConfig)
	if err != nil {
		log.Fatalf("cannot read config: %s", pathToConfig)
	}

	if stat.IsDir() {
		// At first, find a file is named "config.yaml"
		conf := filepath.Join(pathToConfig, "config.yaml")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
		// At second, find a file named "config.default.yaml"
		conf = filepath.Join(pathToConfig, "config.default.yaml")
		if _, err := os.Stat(conf); err == nil {
			return conf
		}
	}

	return pathToConfig
}
