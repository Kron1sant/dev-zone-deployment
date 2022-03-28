package restful

import (
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/sessions"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func GetAppUsers() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, sessions.New(ctx).GetAppUsers())
	}
}

func PostAppUserAction() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		// action from URI
		action := ctx.Param("action")
		var operation func(*dom.User) error
		switch action {
		case "add":
			operation = sessions.New(ctx).AddAppUser
		case "edit":
			operation = sessions.New(ctx).EditAppUser
		case "delete":
			operation = sessions.New(ctx).DeleteAppUser
		case "setPassword":
			SetAppUserPassword(ctx)
			return
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "action can be add or edit",
			})
			return
		}
		// Unmarshaling JSON data from request
		userData := &dom.User{}
		if err := ctx.ShouldBindJSON(userData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Checking app user object
		if err := validateUser(userData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Execute operation
		if err := operation(userData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusCreated, userData)
		}
	}
}

func SetAppUserPassword(ctx *gin.Context) {
	userData := &struct {
		Id       uint   `json:"id"`
		Password string `json:"password"`
	}{}

	// Unmarshaling JSON data from request
	if err := ctx.ShouldBindJSON(userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Execute operation
	if err := sessions.New(ctx).SetAppUserPassword(userData.Id, userData.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.Status(http.StatusCreated)
	}
}

func validateUser(accUser *dom.User) error {
	// check username
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !r.MatchString(accUser.Username) {
		return fmt.Errorf("app username isn't correct: %s", accUser.Username)
	}
	// check email
	r = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if !r.MatchString(accUser.EMail) {
		return fmt.Errorf("e-mail isn't correct: %s", accUser.EMail)
	}
	return nil
}
