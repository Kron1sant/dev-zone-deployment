package config

import (
	"log"

	"github.com/spf13/viper"
)

type dbConfig struct {
	// DataBase host
	Host string
	// DataBase port
	Port string
	// DataBase user to connect
	User string
	// DataBase user's password
	Pass string
	// DataBase name - where the application data is stored
	BaseName string
}

var dbcfg dbConfig

func dbConfInit() {
	baseName := viper.GetString("db.name")
	if baseName == "" {
		baseName = "dev_zone" // default base name
		log.Printf("Database name not specified. Used the default value: %s", baseName)
	}

	dbPass := ""
	if viper.IsSet("DB_PASSWORD") {
		dbPass = viper.GetString("DB_PASSWORD")
	} else {
		dbPass = viper.GetString("db.pass")
	}

	dbcfg = dbConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Pass:     dbPass,
		BaseName: baseName,
	}
}

func GetDBHost() string {
	return dbcfg.Host
}

func GetDBPort() string {
	return dbcfg.Port
}

func GetDBUser() string {
	return dbcfg.User
}

func GetDBPass() string {
	return dbcfg.Pass
}

func GetDBName() string {
	return dbcfg.BaseName
}
