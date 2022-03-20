package api

import (
	"devZoneDeployment/db/dom"
	"fmt"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func IdentityFromUser(user *dom.User) dom.UserIdentity {
	identity := dom.UserIdentity{
		Id:       user.Id,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
	}

	return identity
}

func IdentityFromContext(ctx *gin.Context) dom.UserIdentity {
	return IdentityFromJWTClaims(jwt.ExtractClaims(ctx))
}

func IdentityFromJWTClaims(claims jwt.MapClaims) dom.UserIdentity {
	// check id existence
	idInClaims, ok := claims["id"]
	if !ok {
		// Bad jwt claims format!
		return dom.UserIdentity{
			Empty: true,
		}
	}

	// try converting to uint
	userid, err := strconv.ParseUint(idInClaims.(string), 10, 64)
	if err != nil {
		// Bad jwt claims format!
		return dom.UserIdentity{Empty: true}
	}

	// get isAdmin flag
	isAdminInClaims, ok := claims["isAdmin"]
	if !ok {
		// Bad jwt claims format!
		return dom.UserIdentity{Empty: true}
	}
	isAdmin := isAdminInClaims.(bool)

	// compose identity
	identity := dom.UserIdentity{
		Id:       uint(userid),
		Username: claims["username"].(string),
		IsAdmin:  isAdmin,
	}

	return identity
}

func UserIdentityToJWTClaims(identity dom.UserIdentity) jwt.MapClaims {
	return jwt.MapClaims{
		"id":       fmt.Sprintf("%d", identity.Id),
		"username": identity.Username,
		"isAdmin":  identity.IsAdmin,
	}
}
