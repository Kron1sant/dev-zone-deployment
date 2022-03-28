package restful

import (
	"bytes"
	"devZoneDeployment/config"
	"devZoneDeployment/db/dom"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/alcortesm/tgz"

	"github.com/stretchr/testify/assert"
)

var TEST_DEV_ACCOUNT = dom.DevAccount{
	Id:          333,
	Username:    "test_acc",
	EMail:       "user@mail.foo",
	Name:        "Ivan",
	Patronomic:  "Petrovich",
	Surname:     "Sidorov",
	HasOVPNCert: false,
	Comment:     "Fooooo Bazzz",
}

func TestGetDevAccounts(t *testing.T) {
	assert := assert.New(t)

	// Add two new test accounts
	testAcc1 := dom.DevAccount{
		Username:    "first_acc",
		EMail:       "first@mail.foo",
		Name:        "Petr",
		Patronomic:  "Ivanovich",
		Surname:     "Kii",
		HasOVPNCert: false,
		Comment:     "",
	}
	proccessAccAction("add", testAcc1)
	testAcc2 := dom.DevAccount{
		Username:    "second_acc",
		EMail:       "second@mail.foo",
		Name:        "Vikhail",
		Patronomic:  "Yurievich",
		Surname:     "Stepanov",
		HasOVPNCert: true,
		Comment:     "",
	}
	proccessAccAction("add", testAcc2)

	// Test getting dev accounts
	rr := httptest.NewRecorder()
	ctx := prepareGinContext(rr)
	GetDevAccounts()(ctx)
	// Read response a body
	list := []dom.DevAccount{}
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		assert.FailNow("The response body must contain the data of the dev accounts", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	// The list consists of 2 new test users
	if assert.Len(list, 2, "Response must contain 2 developer accounts") {
		assert.Equal(list[1].Username, testAcc2.Username, "Received Developer account's name must be equal to the source")
		assert.Equal(list[1].EMail, testAcc2.EMail, "Received Developer account's e-mail must be equal to the source")
		assert.Equal(list[1].Name, testAcc2.Name, "Received Developer account's Name must be equal to the source")
		assert.Equal(list[1].Patronomic, testAcc2.Patronomic, "Received Developer account's Patronomic must be equal to the source")
		assert.Equal(list[1].Surname, testAcc2.Surname, "Received Developer account's Surname must be equal to the source")
		assert.Equal(list[1].Comment, testAcc2.Comment, "Received Developer account's Comment must be equal to the source")
		assert.Equal(list[1].HasOVPNCert, testAcc2.HasOVPNCert, "Received Developer account's HasOVPNCert must be equal to the source")
	}
}

func TestPostDevAccountAction(t *testing.T) {
	assert := assert.New(t)
	// The order of the next steps is important
	// Test adding a failed developer account
	stepAddFailedAcc(assert)
	// Test adding a developer account
	testAcc := stepAddCorrectAcc(assert)
	// Test editing an absent developer account
	stepEditFailedAcc(assert, testAcc)
	// Test editing a developer account
	stepEditCorrectAcc(assert, testAcc)
	// Test deleting an absent developer account
	stepDeleteFailedAcc(assert, testAcc)
	// Test deleting a developer account
	stepDeleteCorrectAcc(assert, testAcc)
	// Test executing a bad action
	stepBadAccAction(assert, testAcc)
}

func stepAddFailedAcc(assert *assert.Assertions) {
	wrongAcc := dom.DevAccount{
		Username: "wrong name",
		EMail:    "123@321",
	}
	rr := proccessAccAction("add", wrongAcc)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepAddCorrectAcc(assert *assert.Assertions) dom.DevAccount {
	rr := proccessAccAction("add", TEST_DEV_ACCOUNT)
	// Read response a body
	createdAcc := dom.DevAccount{}
	if err := json.Unmarshal(rr.Body.Bytes(), &createdAcc); err != nil {
		assert.FailNow("The response body must contain the data of the created account", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(TEST_DEV_ACCOUNT.Username, createdAcc.Username, "Created Developer account's name must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.EMail, createdAcc.EMail, "Created Developer account's e-mail must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.Name, createdAcc.Name, "Created Developer account's Name must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.Patronomic, createdAcc.Patronomic, "Created Developer account's Patronomic must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.Surname, createdAcc.Surname, "Created Developer account's Surname must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.Comment, createdAcc.Comment, "Created Developer account's Comment must be equal to the source")
	assert.Equal(TEST_DEV_ACCOUNT.HasOVPNCert, createdAcc.HasOVPNCert, "Created Developer account's HasOVPNCert must be equal to the source")

	return createdAcc
}

func stepEditFailedAcc(assert *assert.Assertions, testAcc dom.DevAccount) {
	wrongAcc := testAcc
	wrongAcc.Id = 999999
	rr := proccessAccAction("edit", wrongAcc)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepEditCorrectAcc(assert *assert.Assertions, testAcc dom.DevAccount) {
	// Change e-mail
	testAcc.EMail = "new@email.baz"
	testAcc.Comment = "New description"
	rr := proccessAccAction("edit", testAcc)
	// Read response a body
	editedAcc := dom.DevAccount{}
	if err := json.Unmarshal(rr.Body.Bytes(), &editedAcc); err != nil {
		assert.FailNow("The response body must contain the data of the edited account", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testAcc.Id, editedAcc.Id, "Edited Developer account's Id must be equal to the source")
	assert.Equal(testAcc.Username, testAcc.Username, "Edited Developer account's name must be equal to the source")
	assert.Equal(testAcc.EMail, testAcc.EMail, "Edited Developer account's e-mail must be equal to the source")
	assert.Equal(testAcc.Name, testAcc.Name, "Edited Developer account's Name must be equal to the source")
	assert.Equal(testAcc.Patronomic, testAcc.Patronomic, "Edited Developer account's Patronomic must be equal to the source")
	assert.Equal(testAcc.Surname, testAcc.Surname, "Edited Developer account's Surname must be equal to the source")
	assert.Equal(testAcc.Comment, testAcc.Comment, "Edited Developer account's Comment must be equal to the source")
	assert.Equal(testAcc.HasOVPNCert, testAcc.HasOVPNCert, "Edited Developer account's HasOVPNCert must be equal to the source")
}

func stepDeleteFailedAcc(assert *assert.Assertions, testAcc dom.DevAccount) {
	wrongAcc := testAcc
	wrongAcc.Id = 999999
	rr := proccessAccAction("delete", wrongAcc)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepDeleteCorrectAcc(assert *assert.Assertions, testAcc dom.DevAccount) {
	rr := proccessAccAction("delete", testAcc)
	// Read response a body
	deletedAcc := dom.DevAccount{}
	if err := json.Unmarshal(rr.Body.Bytes(), &deletedAcc); err != nil {
		assert.FailNow("The response body must contain the data of the deleted user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testAcc.Id, deletedAcc.Id, "Deleted Developer account's Id must be equal to the source")
	assert.Equal(testAcc.Username, deletedAcc.Username, "Deleted Developer account's name must be equal to the source")
	assert.Equal(testAcc.EMail, deletedAcc.EMail, "Deleted Developer account's e-mail must be equal to the source")
}

func stepBadAccAction(assert *assert.Assertions, testApp dom.DevAccount) {
	rr := proccessAccAction("bad_action", testApp)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func proccessAccAction(action string, acc dom.DevAccount) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	ctx := prepareGinContext(w)
	setAction(ctx, action)
	b, _ := json.Marshal(acc)
	ctx.Request = httptest.NewRequest("POST", "/accounts/"+action, bytes.NewReader(b))
	PostDevAccountAction()(ctx)
	return w
}

func TestPostOpenVPNKey(t *testing.T) {
	assert := assert.New(t)

	// Deploy a temporary environment for easy-rsa
	workspace, err := deployTestEasyRSAWorkspace()
	if workspace != "" {
		defer os.RemoveAll(workspace)
	}
	if err != nil {
		assert.FailNow("Failed to deploy easy-rsa", err)
	}

	// Add the test dev account that will getting the openvpn key
	rr := proccessAccAction("add", TEST_DEV_ACCOUNT)
	// Read response a body
	createdAcc := dom.DevAccount{}
	if err := json.Unmarshal(rr.Body.Bytes(), &createdAcc); err != nil {
		assert.FailNow("The response body must contain the data of the created account", err)
	}

	// Test getting dev accounts
	rr = httptest.NewRecorder()
	ctx := prepareGinContext(rr)
	b, _ := json.Marshal(createdAcc)
	ctx.Request = httptest.NewRequest("POST", "/openvpnkey", bytes.NewReader(b))
	PostOpenVPNKey()(ctx)
	// Read response a body
	ovpnKey := []byte{}
	if err := json.Unmarshal(rr.Body.Bytes(), &ovpnKey); err != nil {
		assert.FailNow("The response body must contain the OpenVPN key", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	assert.Equal(bytes.Index(ovpnKey, []byte("client")), 0, "OpenVPN key must be started with 'client'")
}

func deployTestEasyRSAWorkspace() (string, error) {
	distr := "https://github.com/OpenVPN/easy-rsa/releases/download/v3.0.8/EasyRSA-3.0.8.tgz"
	distrpath, err := ioutil.TempFile("", "EasyRSA.tgz")
	if err != nil {
		return "", err
	}
	defer os.Remove(distrpath.Name()) // clean up
	if err := downloadEasyRSA(distr, distrpath.Name()); err != nil {
		return "", err
	}
	tmpPath, err := tgz.Extract(distrpath.Name())
	if err != nil {
		return "", err
	}

	// Initialize workspace
	ws := path.Join(tmpPath, "EasyRSA-3.0.8")
	ca := path.Join(ws, "pki", "ca.crt")
	tls := path.Join(ws, "tls.key")
	template := path.Join(ws, "client-common.txt")
	cmd(ws, "./easyrsa", "init-pki")
	cmd(ws, "./easyrsa", "--batch", "build-ca", "nopass")
	cmd(ws, "./easyrsa", "build-server-full", "server", "nopass")
	cmd(ws, "./easyrsa", "gen-crl")

	f_tls, err := os.Create(tls)
	if err != nil {
		return "", err
	}
	// Mocking data
	fmt.Fprintf(f_tls, "-----BEGIN OpenVPN Static key V1-----\n")
	fmt.Fprintf(f_tls, "FooooBazzzz\n")
	fmt.Fprintf(f_tls, "-----END OpenVPN Static key V1-----\n")

	f_tmplt, err := os.Create(template)
	if err != nil {
		return "", err
	}
	fmt.Fprintf(f_tmplt, "client\n")

	// Set config
	ovpnParams := fmt.Sprintf(`{
		"easyrsa": "%s",
		"ca-cert": "%s",
		"tls-key": "%s",
		"client-template": "%s"
	}`, ws, ca, tls, template)
	config.SetOpenVPNParamsFromJSON(ovpnParams)

	return ws, nil
}

func downloadEasyRSA(url, distrpath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(distrpath)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	return nil
}

func cmd(ws string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = ws
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
