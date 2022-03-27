package restful

import (
	"bytes"
	"devZoneDeployment/api"
	"devZoneDeployment/db/dom"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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
	w := httptest.NewRecorder()
	ctx := prepareContext(w)
	GetAppUsers()(ctx)
	fmt.Println(w.Body.String())
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
		Username: "worng name",
		EMail:    "123@321",
	}
	rr := proccessAction("add", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepAddCorrectUser(assert *assert.Assertions) dom.User {
	rr := proccessAction("add", TEST_APP_USER)
	// Read response a body
	createdUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &createdUser); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(TEST_APP_USER.Username, createdUser.Username, "User name must be equal")
	assert.Equal(TEST_APP_USER.EMail, createdUser.EMail, "User e-mail must be equal")
	assert.Equal(createdUser.HasDevAccount, createdUser.HasDevAccount, "User HasDevAccount must be equal")
	assert.Equal(createdUser.DevAccountId, createdUser.DevAccountId, "User DevAccountId must be equal")

	return createdUser
}

func stepEditFailedUser(assert *assert.Assertions, testUser dom.User) {
	wrongUser := testUser
	wrongUser.Id = 999999
	rr := proccessAction("edit", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepEditCorrectUser(assert *assert.Assertions, testUser dom.User) {
	// Change e-mail
	testUser.EMail = "new@email.baz"
	rr := proccessAction("edit", testUser)
	// Read response a body
	editedUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &editedUser); err != nil {
		assert.FailNow("The response body must contain the data of the edited user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testUser.Id, editedUser.Id, "User Id must be equal")
	assert.Equal(testUser.Username, editedUser.Username, "User name must be equal")
	assert.Equal(testUser.EMail, editedUser.EMail, "User e-mail must be equal")
	assert.Equal(testUser.HasDevAccount, editedUser.HasDevAccount, "User HasDevAccount must be equal")
	assert.Equal(testUser.DevAccountId, editedUser.DevAccountId, "User DevAccountId must be equal")
}

func stepSetupPassword(assert *assert.Assertions, testUser dom.User) {
	rr := setUserPassword(testUser, "new_pass")
	// Check the code
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
}

func stepDeleteFailedUser(assert *assert.Assertions, testUser dom.User) {
	wrongUser := testUser
	wrongUser.Id = 999999
	rr := proccessAction("delete", wrongUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepDeleteCorrectUser(assert *assert.Assertions, testUser dom.User) {
	rr := proccessAction("delete", testUser)
	// Read response a body
	deletedUser := dom.User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &deletedUser); err != nil {
		assert.FailNow("The response body must contain the data of the deleted user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testUser.Id, deletedUser.Id, "User Id must be equal")
	assert.Equal(testUser.Username, deletedUser.Username, "User name must be equal")
	assert.Equal(testUser.EMail, deletedUser.EMail, "User e-mail must be equal")
}

func stepBadAction(assert *assert.Assertions, testUser dom.User) {
	rr := proccessAction("bad_action", testUser)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the data of the created user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func proccessAction(action string, user dom.User) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	ctx := prepareContext(w)
	setAction(ctx, action)
	b, _ := json.Marshal(user)
	ctx.Request = httptest.NewRequest("POST", "/accounts/"+action, bytes.NewReader(b))
	PostAppUserAction()(ctx)
	return w
}

func prepareContext(w http.ResponseWriter) *gin.Context {
	ctx, _ := gin.CreateTestContext(w)
	uid := api.UserIdentity{
		Id:       123,
		Username: TEST_APP_ADMIN_NAME,
		IsAdmin:  true,
	}
	jwtClaims := api.UserIdentityToJWTClaims(uid)
	ctx.Set("JWT_PAYLOAD", jwtClaims)
	return ctx
}

func setAction(ctx *gin.Context, action string) {
	ctx.Params = []gin.Param{
		{
			Key:   "action",
			Value: action,
		},
	}
}

func setUserPassword(user dom.User, pass string) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	ctx := prepareContext(rr)
	passwdData := fmt.Sprintf(`{"id": %d, "password": "%s"}`, user.Id, pass)
	ctx.Request = httptest.NewRequest("POST", "/accounts/setPassword", strings.NewReader(passwdData))
	setAction(ctx, "setPassword")
	PostAppUserAction()(ctx)
	return rr
}
