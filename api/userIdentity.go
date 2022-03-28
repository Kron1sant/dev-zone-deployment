package api

import (
	"devZoneDeployment/db/dom"
	"fmt"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// Represent an identity of an app user session
// Using to authorize actions in a session
type UserIdentity struct {
	Id       uint
	Username string
	IsAdmin  bool
	Empty    bool
	Guest    bool
}

// Getting UserIdentity by User
func IdentityFromUser(user *dom.User) UserIdentity {
	identity := UserIdentity{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	return identity
}

// Getting UserIdentity from GIN context with JWT
func IdentityFromContext(ctx *gin.Context) UserIdentity {
	return IdentityFromJWTClaims(jwt.ExtractClaims(ctx))
}

// Getting UserIdentity from JWT claims
func IdentityFromJWTClaims(claims jwt.MapClaims) UserIdentity {
	// check id existence
	idInClaims, ok := claims["id"]
	if !ok {
		// Bad jwt claims format!
		return UserIdentity{
			Empty: true,
		}
	}

	// try converting to uint
	userid, err := strconv.ParseUint(idInClaims.(string), 10, 64)
	if err != nil {
		// Bad jwt claims format!
		return UserIdentity{Empty: true}
	}

	// get isAdmin flag
	isAdminInClaims, ok := claims["isAdmin"]
	if !ok {
		// Bad jwt claims format!
		return UserIdentity{Empty: true}
	}
	isAdmin := isAdminInClaims.(bool)

	// compose identity
	identity := UserIdentity{
		Id:       uint(userid),
		Username: claims["username"].(string),
		IsAdmin:  isAdmin,
	}

	return identity
}

// Preparing JWT claims with UserIdentity
func UserIdentityToJWTClaims(identity UserIdentity) jwt.MapClaims {
	return jwt.MapClaims{
		"id":       fmt.Sprintf("%d", identity.Id),
		"username": identity.Username,
		"isAdmin":  identity.IsAdmin,
	}
}
