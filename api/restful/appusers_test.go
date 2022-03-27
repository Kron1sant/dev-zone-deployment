package restful

import (
	"bytes"
	"devZoneDeployment/db/dom"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TEST_APP_USER = dom.User{
	Id:       333,
	Username: "test_user",
	EMail:    "user@mail.foo",
	Password: "test_user_pass",
	IsAdmin:  true,
}

func TestGetAppUsers(t *testing.T) {
	assert := assert.New(t)

	// Add two new test users
	testUser1 := dom.User{
		Username: "first_user",
		EMail:    "foo@mail.foo",
		IsAdmin:  true,
	}
	proccessUserAction("add", testUser1)
	testUser2 := dom.User{
		Username: "second_user",
		EMail:    "baz@mail.baz",
		IsAdmin:  true,
	}
	proccessUserAction("add", testUser2)

	// Test getting app users
	rr := httptest.NewRecorder()
	ctx := prepareGinContext(rr)
	GetAppUsers()(ctx)
	// Read response a body
	list := []dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &list); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	// The list consists of admin + 2 new test users = 3
	if assert.Len(list, 3, "Response must containe 3 users") {
		assert.Equal(list[2].Username, testUser2.Username, "User name must be equal to the source")
		assert.Equal(list[2].EMail, testUser2.EMail, "User e-mail must be equal to the source")
		assert.Equal(list[2].IsAdmin, testUser2.IsAdmin, "User IsAdmin must be equal to the source")
		assert.Equal(list[2].HasDevAccount, testUser2.HasDevAccount, "User HasDevAccount must be equal to the source")
		assert.Equal(list[2].DevAccountId, testUser2.DevAccountId, "User DevAccountId must be equal to the source")
	}
}

func TestPostAppUserAction(t *testing.T) {
	assert := assert.New(t)
	// The order of the next steps is important
	// Test adding a failed user
	stepAddFailedUser(assert)
	// Test adding an user
	testUser := stepAddCorrectUser(assert)
	// Test editing an absent user
	stepEditFailedUser(assert, testUser)
	// Test editing an user
	stepEditCorrectUser(assert, testUser)
	// Test seting an user's password
	stepSetupPassword(assert, testUser)
	// Test deleting an absent user
	stepDeleteFailedUser(assert, testUser)
	// Test deleting an user
	stepDeleteCorrectUser(assert, testUser)
	// Test executing a bad action
	stepBadAction(assert, testUser)
}

func stepAddFailedUser(assert *assert.Assertions) {
	wrongUser := dom.User{
		Username: "wrong name",
		EMail:    "123@321",
	}
	rr := proccessUserAction("add", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepAddCorrectUser(assert *assert.Assertions) dom.User {
	rr := proccessUserAction("add", TEST_APP_USER)
	// Read response a body
	createdUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &createdUser); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(TEST_APP_USER.Username, createdUser.Username, "Created User's name must be equal to the source")
	assert.Equal(TEST_APP_USER.EMail, createdUser.EMail, "Created User's e-mail must be equal to the source")
	assert.Equal(TEST_APP_USER.IsAdmin, createdUser.IsAdmin, "Created User's IsAdmin must be equal to the source")
	assert.Equal(TEST_APP_USER.HasDevAccount, createdUser.HasDevAccount, "Created User's HasDevAccount must be equal to the source")
	assert.Equal(TEST_APP_USER.DevAccountId, createdUser.DevAccountId, "Created User's DevAccountId must be equal to the source")

	return createdUser
}

func stepEditFailedUser(assert *assert.Assertions, testUser dom.User) {
	wrongUser := testUser
	wrongUser.Id = 999999
	rr := proccessUserAction("edit", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepEditCorrectUser(assert *assert.Assertions, testUser dom.User) {
	// Change e-mail
	testUser.EMail = "new@email.baz"
	rr := proccessUserAction("edit", testUser)
	// Read response a body
	editedUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &editedUser); err != nil {
		assert.FailNow("The response body must contain the data of the edited user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testUser.Id, editedUser.Id, "Edited User's Id must be equal to the source")
	assert.Equal(testUser.Username, editedUser.Username, "Edited User's name must be equal to the source")
	assert.Equal(testUser.EMail, editedUser.EMail, "Edited User's e-mail must be equal to the source")
	assert.Equal(testUser.HasDevAccount, editedUser.HasDevAccount, "Edited User's HasDevAccount must be equal to the source")
	assert.Equal(testUser.DevAccountId, editedUser.DevAccountId, "Edited User's DevAccountId must be equal to the source")
}

func stepSetupPassword(assert *assert.Assertions, testUser dom.User) {
	rr := setUserPassword(testUser, "new_pass")
	// Check the code
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
}

func stepDeleteFailedUser(assert *assert.Assertions, testUser dom.User) {
	wrongUser := testUser
	wrongUser.Id = 999999
	rr := proccessUserAction("delete", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepDeleteCorrectUser(assert *assert.Assertions, testUser dom.User) {
	rr := proccessUserAction("delete", testUser)
	// Read response a body
	deletedUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &deletedUser); err != nil {
		assert.FailNow("The response body must contain the data of the deleted user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testUser.Id, deletedUser.Id, "Deleted User's Id must be equal to the source")
	assert.Equal(testUser.Username, deletedUser.Username, "Deleted User's name must be equal to the source")
	assert.Equal(testUser.EMail, deletedUser.EMail, "Deleted User's e-mail must be equal to the source")
}

func stepBadAction(assert *assert.Assertions, testUser dom.User) {
	rr := proccessUserAction("bad_action", testUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func proccessUserAction(action string, user dom.User) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	ctx := prepareGinContext(w)
	setAction(ctx, action)
	b, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest("POST", "/users/"+action, bytes.NewReader(b))
	PostAppUserAction()(ctx)
	return w
}

func setUserPassword(user dom.User, pass string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ctx := prepareGinContext(rr)
	passwdData := fmt.Sprintf(`{"id": %d, "password": "%s"}`, user.Id, pass)
	ctx.Request = httptest.NewRequest("POST", "/accounts/setPassword", strings.NewReader(passwdData))
	setAction(ctx, "setPassword")
	PostAppUserAction()(ctx)
	return rr
}
