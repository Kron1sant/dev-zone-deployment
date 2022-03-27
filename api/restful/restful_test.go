package restful

import (
	"devZoneDeployment/api"
	"devZoneDeployment/config"
	"devZoneDeployment/db/mongodb"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

const (
	TEST_DB_NAME        = "test_dev_zone"
	TEST_APP_ADMIN_NAME = "test_admin"
	TEST_APP_ADMIN_PASS = "test_pass_123"
)

func TestMain(m *testing.M) {
	setup()         // before testing
	code := m.Run() // starts testing
	teardown()      // after testing
	os.Exit(code)
}

func setup() {
	dbparams := `{"host":"localhost",
		"port":"27017",
		"user":"user",
		"pass":"pass",
		"basename":"` + TEST_DB_NAME + `"
	}`
	appparams := `{"port":"8089",
		"secret":"test_secret_key",
		"default_admin": {
			"username":"` + TEST_APP_ADMIN_NAME + `",
			"password":"` + TEST_APP_ADMIN_PASS + `"
		}
	}`
	config.SetDBParamsFromJSON(dbparams)
	config.SetAppParamsFromJSON(appparams)

	// During initialization, a new test base will be created
	// and the specified user will be started
	mongodb.UseMongoDBSource()
}

func teardown() {
	// Delete test database
	if mongodb.DS.Database.Name() == TEST_DB_NAME {
		mongodb.DS.Database.Drop(mongodb.DefaulContext())
	} else {
		log.Fatalf("the name of the database differs from the expected value: %s != %s",
			mongodb.DS.Database.Name(), TEST_DB_NAME)
	}
}

// request sends req to the mockRouter and returns Response
func request(mockRouter *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mockRouter.ServeHTTP(rr, req)
	return rr
}

// setAction sets context params with specified action
func setAction(ctx *gin.Context, action string) {
	ctx.Params = []gin.Param{
		{
			Key:   "action",
			Value: action,
		},
	}
}

// prepareGinContext creates a test gin context
func prepareGinContext(w http.ResponseWriter) *gin.Context {
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
