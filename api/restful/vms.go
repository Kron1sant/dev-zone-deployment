package restful

import (
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/sessions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListVM() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, sessions.New(ctx).ListVirtualMachines())
	}
}

func PostVMAction() func(*gin.Context) {
	// Using closure to passing datasource in handler
	return func(ctx *gin.Context) {
		// action from URI
		action := ctx.Param("action")
		var operation func(*dom.VM) error
		switch action {
		case "update":
			UpdateListVirtualMachinesFromCloud(ctx)
			return
		case "":
			// ToDo
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "action can be add or edit",
			})
			return
		}
		// Unmarshaling JSON data from request
		vmData := &dom.VM{}
		if err := ctx.ShouldBindJSON(vmData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Checking app user object
		if err := validateVM(vmData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Execute operation
		if err := operation(vmData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			ctx.JSON(http.StatusCreated, vmData)
		}
	}
}

func UpdateListVirtualMachinesFromCloud(ctx *gin.Context) {
	// Execute operation
	if err := sessions.New(ctx).UpdateListVirtualMachinesFromCloud(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.Status(http.StatusCreated)
	}
}

func validateVM(vm *dom.VM) error {
	// check vm - ToDo
	return nil
}
