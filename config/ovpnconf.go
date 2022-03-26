package config

import "github.com/spf13/viper"

type ovpnConfig struct {
	// Path to an easyrsa directory with bin data
	Easyrsa string
	// Path to ca cert
	CaCert string
	// Path to tls key
	TlsKey string
	// Path to file with the client key template
	ClientTemplate string
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
