package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

type ovpnConfig struct {
	// Path to an easyrsa directory with bin data
	Easyrsa string `json:"easyrsa"`
	// Path to ca cert
	CaCert string `json:"ca-cert"`
	// Path to tls key
	TlsKey string `json:"tls-key"`
	// Path to file with the client key template
	ClientTemplate string `json:"client-template"`
}

var ovpncfg ovpnConfig

func ovpnConfInit() {
	ovpncfg = ovpnConfig{
		Easyrsa:        viper.GetString("openvpn.easyrsa"),
		CaCert:         viper.GetString("openvpn.ca-cert"),
		TlsKey:         viper.GetString("openvpn.tls-key"),
		ClientTemplate: viper.GetString("openvpn.client-template"),
	}
}

func GetEasyRsaDir() string {
	return ovpncfg.Easyrsa
}

func GetCaCertPath() string {
	return ovpncfg.CaCert
}

func GetTlsKeyPath() string {
	return ovpncfg.TlsKey
}

func GetClientTemplatePath() string {
	return ovpncfg.ClientTemplate
}

func SetOpenVPNParamsFromJSON(params string) {
	if err := json.Unmarshal([]byte(params), &ovpncfg); err != nil {
		log.Fatalf("set openvpn params from json failed: %s", err)
	}
}
