package restful

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAddAuthMiddleware(t *testing.T) {
	mockRouter := gin.New()
	AddAuthMiddleware(mockRouter)
	// Tests /login
	t.Run("Success login test", func(t *testing.T) { successLoginTest(t, mockRouter) })
	t.Run("Failed login test", func(t *testing.T) { failedLoginTest(t, mockRouter) })
	// Tests /logout
	t.Run(" Logout test", func(t *testing.T) { logoutTest(t, mockRouter) })
	// Tests /auth/refresh_token
	t.Run("Success refresh token test", func(t *testing.T) { successRefreshTokenTest(t, mockRouter) })
	t.Run("Failed refresh token test", func(t *testing.T) { failedRefreshTokenTest(t, mockRouter) })
	// Tests /auth/check
	t.Run("Success authentification check test", func(t *testing.T) { successAuthCheckTest(t, mockRouter) })
	t.Run("Failed authentification check test", func(t *testing.T) { failedAuthCheckTest(t, mockRouter) })
}

// Test case "Successful login"
func successLoginTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	rr, err := tryLoginToApp(mockRouter, TEST_APP_ADMIN_NAME, TEST_APP_ADMIN_PASS)
	assert.NoError(err, "test failed: %s", err)
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	assert.NotEmpty(rr.Header().Get("Set-Cookie"), "Cookie must be set")
}

// Test case "Failed login"
func failedLoginTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	wrongPass := "wrong pass"
	rr, err := tryLoginToApp(mockRouter, TEST_APP_ADMIN_NAME, wrongPass)
	assert.NoError(err, "test failed: %s", err)
	assert.Equal(http.StatusUnauthorized, rr.Code, "The response code must be 401")
}

// Test case "Logout"
func logoutTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	res, err := tryLoginToApp(mockRouter, TEST_APP_ADMIN_NAME, TEST_APP_ADMIN_PASS)
	if err != nil || res.Code != http.StatusOK {
		// Failed login. Do nothing, because there is another according test
		return
	}

	credCookie := res.Header().Get("Set-Cookie")
	req := httptest.NewRequest("POST", "/logout", nil)
	req.Header.Add("Cookie", credCookie)

	rr := request(mockRouter, req)
	assert.Equal(rr.Code, http.StatusOK, "The response code must be 200")
	assert.Contains(rr.Header().Get("Set-Cookie"), "token=;", "After Logout cookie must have 'token=;'")
}

// Test case "Successful refresh token"
func successRefreshTokenTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	res, err := tryLoginToApp(mockRouter, TEST_APP_ADMIN_NAME, TEST_APP_ADMIN_PASS)
	if err != nil || res.Code != http.StatusOK {
		// Failed login. Do nothing, because there is another according test
		return
	}

	credCookie := res.Header().Get("Set-Cookie")
	req := httptest.NewRequest("GET", "/auth/refresh_token", nil)
	req.Header.Add("Cookie", credCookie)

	rr := request(mockRouter, req)
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	assert.NotEmpty(rr.Header().Get("Set-Cookie"), "Cookie must be set")
}

// Test case "Failed refresh token"
func failedRefreshTokenTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	req := httptest.NewRequest("GET", "/auth/refresh_token", nil)
	req.Header.Add("Cookie", "token=wrong_token;")

	rr := request(mockRouter, req)
	assert.Equal(http.StatusUnauthorized, rr.Code, "The response code must be 401")
}

// Test case "Successful checking authentication"
func successAuthCheckTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	res, err := tryLoginToApp(mockRouter, TEST_APP_ADMIN_NAME, TEST_APP_ADMIN_PASS)
	if err != nil || res.Code != http.StatusOK {
		// Failed login. Do nothing, because there is another according test
		return
	}

	credCookie := res.Header().Get("Set-Cookie")
	req := httptest.NewRequest("POST", "/auth/check", nil)
	req.Header.Add("Cookie", credCookie)

	rr := request(mockRouter, req)
	payload := struct {
		IsAdmin  bool   `json:"isAdmin"`
		Username string `json:"username"`
	}{}
	assert.NoError(json.Unmarshal(rr.Body.Bytes(), &payload), "Failed unmarshalling payload from jwt")
	assert.Equal(http.StatusOK, rr.Code, "The response code must be 200")
	assert.True(payload.IsAdmin, "The response must content isAdmin flag")
	assert.Equal(TEST_APP_ADMIN_NAME, payload.Username, "The response code must content username")
}

// Test case "Failed checking authentication"
func failedAuthCheckTest(t *testing.T, mockRouter *gin.Engine) {
	assert := assert.New(t)
	req := httptest.NewRequest("POST", "/auth/check", nil)
	req.Header.Add("Cookie", "token=wrong_token;")

	rr := request(mockRouter, req)
	assert.Equal(http.StatusUnauthorized, rr.Code, "The response code must be 401")
}

// tryLoginToApp emulates a user authorization attempt.
// Returns a response and an error
func tryLoginToApp(mockRouter *gin.Engine, user string, pass string) (*httptest.ResponseRecorder, error) {
	cred, err := json.Marshal(credentialData{
		Username: user,
		Password: pass,
	})
	if err != nil {
		return nil, fmt.Errorf("marshalling credentials failed: %s", err)
	}

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(cred))
	req.Header.Add("Content-Type", "application/json")

	return request(mockRouter, req), nil
}
