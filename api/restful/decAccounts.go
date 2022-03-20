package restful

import (
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/sessions"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func GetDevAccounts() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, sessions.New(ctx).GetDevAccounts())
	}
}

func PostDevAccountAction() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		// action from URI
		action := ctx.Param("action")
		var operation func(*dom.DevAccount) error
		switch action {
		case "add":
			operation = sessions.New(ctx).AddDevAccount
		case "edit":
			operation = sessions.New(ctx).EditDevAccount
		case "delete":
			operation = sessions.New(ctx).DeleteDevAccount
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "action can be wether add or edit",
			})
			return
		}
		// Unmarshaling JSON data from request
		newAccount := &dom.DevAccount{}
		if err := ctx.ShouldBindJSON(newAccount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Checking devaccount object
		if err := validateAccount(newAccount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Execute operation
		if err := operation(newAccount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusCreated, newAccount)
		}
	}
}

func PostOpenVPNKey() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		// Unmarshaling JSON data from request
		newAccount := &dom.DevAccount{}
		if err := ctx.ShouldBindJSON(newAccount); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		data, err := sessions.New(ctx).GetNewOrExistsOpenVPNKey(newAccount)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, data)
		}
	}
}

func validateAccount(devAccount *dom.DevAccount) error {
	// check username
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	if !r.MatchString(devAccount.Username) {
		return fmt.Errorf("dev username isn't correct: %s", devAccount.Username)
	}
	// check email
	r = regexp.MustCompile(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	if !r.MatchString(devAccount.EMail) {
		return fmt.Errorf("e-mail isn't correct: %s", devAccount.EMail)
	}
	return nil
}
