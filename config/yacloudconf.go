package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

type yacloudConfig struct {
	// Inner id of the folder of the cloud
	Folderid string `json:"folderId"`
	// Yandex cloud token
	Token string `json:"token"`
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

func SetYaCloudParamsFromJSON(params string) {
	if err := json.Unmarshal([]byte(params), &yacfg); err != nil {
		log.Fatalf("set db params from json failed: %s", err)
	}
}
