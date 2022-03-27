package restful

import (
	"bytes"
	"devZoneDeployment/db/dom"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
