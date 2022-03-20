package devworkspace

import (
	"bufio"
	"bytes"
	"devZoneDeployment/config"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

func GetOpenVPNKey(keyName string) ([]byte, error) {
	if vpnKeyExists(keyName) {
		return composeOpenVPNFile(keyName)
	} else {
		return generateNewOpenVPNKey(keyName)
	}
}

func vpnKeyExists(keyName string) bool {
	path2index := path.Join(config.AppConfig.GetString("openvpn.easyrsa"), "pki", "index.txt")

	f, err := os.Open(path2index)
	if err != nil {
		log.Printf("cannot open user-index file: %s: %s", path2index, err.Error())
		return false
	}
	defer f.Close()

	scaner := bufio.NewScanner(f)
	for scaner.Scan() {
		if strings.Contains(scaner.Text(), "/CN="+keyName) {
			return true
		}
	}

	return false
}

func generateNewOpenVPNKey(keyName string) ([]byte, error) {
	easyrsa := path.Join(config.AppConfig.GetString("openvpn.easyrsa"), "easyrsa")
	cmd := exec.Command(easyrsa, "build-client-full", keyName, "nopass")
	cmd.Dir = config.AppConfig.GetString("openvpn.easyrsa")

	if err := cmd.Run(); err != nil {
		log.Printf("cannot generate key: %s", cmd.String())
		return nil, err
	}

	// after generation key compose all cert and private files to *.ovpn
	return composeOpenVPNFile(keyName)
}

func composeOpenVPNFile(keyName string) ([]byte, error) {
	templatePath := config.AppConfig.GetString("openvpn.client-template")
	template, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Printf("cannot compose openvpn file: %s: %s", templatePath, err.Error())
		return nil, err
	}

	// collect all parts of an ovpn key in the "res"
	res := bytes.NewBuffer(template)

	fmt.Fprintln(res, "<ca>")
	caPath := config.AppConfig.GetString("openvpn.ca-cert")
	ca, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Printf("cannot compose openvpn file (ca certificate): %s: %s", caPath, err.Error())
		return nil, err
	}
	res.Write(ca)
	fmt.Fprintln(res, "</ca>")

	easyrsaPath := config.AppConfig.GetString("openvpn.easyrsa")

	fmt.Fprintln(res, "<cert>")
	userCertPath := path.Join(easyrsaPath, "pki", "issued", keyName+".crt")
	data, err := readPayloadFromFile(userCertPath, "BEGIN CERTIFICATE")
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(res, data)
	fmt.Fprintln(res, "</cert>")

	fmt.Fprintln(res, "<key>")
	userKeyPath := path.Join(easyrsaPath, "pki", "private", keyName+".key")
	key, err := ioutil.ReadFile(userKeyPath)
	if err != nil {
		log.Printf("cannot compose openvpn file (user certificate): %s: %s", userKeyPath, err.Error())
		return nil, err
	}
	res.Write(key)
	fmt.Fprintln(res, "</key>")

	fmt.Fprintln(res, "<tls-crypt>")
	tlsKey := config.AppConfig.GetString("openvpn.tls-key")
	data, err = readPayloadFromFile(tlsKey, "BEGIN OpenVPN Static key")
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(res, data)
	fmt.Fprint(res, "</tls-crypt>")

	return res.Bytes(), nil
}

func readPayloadFromFile(filepath string, startArea string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Printf("cannot compose openvpn file: %s: %s", filepath, err.Error())
		return "", err
	}
	var res strings.Builder
	scanner := bufio.NewScanner(f)
	payload := false
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), startArea) {
			payload = true
		}
		if payload {
			res.WriteString(scanner.Text() + "\n")
		}
	}

	return res.String(), nil
}
