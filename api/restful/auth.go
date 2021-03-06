package restful

import (
	"devZoneDeployment/api"
	"devZoneDeployment/config"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/sessions"
	"log"
	"math/rand"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Use for loggin in user
type credentialData struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

// Auth middleware controls access to restricted service methods by JWT token.
// Carries out the release and validate JWT.
// Adds handlers to followed URIs:
// 	* /login - logging in by username and pass. Return JWT by JSON response and cookie
//  * /logout - logging off, forgot token
//	* /auth/refresh_token - reissue expiring token
//	* /auth/check - check the freshnest of the token
// Return IRouter that serves /auth path
func AddAuthMiddleware(router *gin.Engine) gin.IRouter {
	authMiddleware, err := jwt.New(prepareJWTMiddleware())
	if err != nil {
		log.Fatalf("JWT error: %s\n", err.Error())
	}
	router.POST("/login", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)

	auth := router.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.POST("/check", checkAuthHandler)
	return auth
}

func prepareJWTMiddleware() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:           "DevZoneDeployer",
		Key:             getSecretKey(),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityHandler: identityHandler,
		PayloadFunc:     payloadFunc,
		Authenticator:   authenticatorHandler,
		Unauthorized:    unauthorizedHandler,
		TokenLookup:     "cookie:token, header:Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		SendCookie:      true,
		SecureCookie:    false, //non HTTPS dev environments
		CookieHTTPOnly:  true,  // JS can't modify
		//CookieDomain:   "localhost:8080",
		CookieName:     "token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
	}
}

func getSecretKey() []byte {
	// Get secret from config
	secret := config.GetAppSecret()
	if secret == "" {
		// If secret is empty gen random secret
		s := make([]byte, 20)
		rand.Seed(time.Now().Unix())
		for i := 0; i < 20; i++ {
			s[i] = byte(rand.Intn(94) + 33)
		}
		secret = string(s)
	}
	return []byte(secret)
}

func identityHandler(c *gin.Context) interface{} {
	// Exrtact claims from JWT and convert to UserIdentity
	return api.IdentityFromContext(c)
}

// Payload will save in a token
func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*dom.User); ok {
		// Pack User to JWT claims by UserIdentity
		return api.UserIdentityToJWTClaims(api.IdentityFromUser(v))
	}
	return jwt.MapClaims{}
}

func authenticatorHandler(ctx *gin.Context) (interface{}, error) {
	var login credentialData
	if err := ctx.ShouldBind(&login); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	user, err := sessions.NewGuest(ctx).Authenticate(login.Username, login.Password)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	} else {
		return user, nil
	}
}

func unauthorizedHandler(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

// checkAuthHandler extracts UserIdentity from current context,
// checks that the corresponding user is in the database,
// and returns a response with user parameters or an error
func checkAuthHandler(ctx *gin.Context) {
	userIdentity := api.IdentityFromContext(ctx)
	if userIdentity.Empty {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad jwt, cannot identify the user",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"username": userIdentity.Username,
			"isAdmin":  userIdentity.IsAdmin,
		})
	}
}
