package mongodb

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/devworkspace"
	"fmt"
	"strings"
)

func (ds *MongoDBSource) GetNewOrExistsOpenVPNKey(uid api.UserIdentity, acc *dom.DevAccount) ([]byte, error) {
	if !uid.IsAdmin && uid.Id != acc.Id {
		return nil, fmt.Errorf("you haven't got access to this dev account")
	}

	if !ds.devAccountExists(uid, acc, false) {
		return nil, fmt.Errorf("the account doesn't exist")
	}

	if !acc.HasOVPNCert {
		// We got a new ovpn key, save this information in DB
		acc.HasOVPNCert = true
		ds.SetDevAccounts(uid, acc, false)
	}

	keyName := acc.OpenVPNKeyName

	if keyName == "" {
		keyName = fmt.Sprintf("%d-%s", acc.Id, strings.ToLower(strings.TrimSpace(acc.Username)))
	}

	res, err := devworkspace.GetOpenVPNKey(keyName)
	if err == nil && keyName != acc.OpenVPNKeyName {
		// if it is the first openVPN key, then save it name after succesful generation
		acc.OpenVPNKeyName = keyName
		ds.SetDevAccounts(uid, acc, false)
	}
	return res, err
}
