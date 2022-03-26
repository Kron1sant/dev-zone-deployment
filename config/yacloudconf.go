package config

import "github.com/spf13/viper"

type yacloudConfig struct {
	// Inner id of the folder of the cloud
	Folderid string
	// Yandex cloud token
	Token string
}

var yacfg yacloudConfig

func yacloudConfInit() {
	token := ""
	if viper.IsSet("DB_PASSWORD") {
		token = viper.GetString("DB_PASSWORD")
	} else {
		token = viper.GetString("yacloud.token")
	}

	yacfg = yacloudConfig{
		Folderid: viper.GetString("yacloud.folderid"),
		Token:    token,
	}
}

func GetYaFolderID() string {
	return yacfg.Folderid
}

func GetYaToken() string {
	return yacfg.Token
}
