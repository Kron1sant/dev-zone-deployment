package restful

import (
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/sessions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVMs() func(*gin.Context) {
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
		case "add":
			operation = sessions.New(ctx).AddVirtualMachine
		case "edit":
			operation = sessions.New(ctx).EditVirtualMachine
		case "delete":
			operation = sessions.New(ctx).DeleteVirtualMachine
		case "update":
			UpdateListVirtualMachinesFromCloud(ctx)
			return
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
	// check status
	availablStatuses := []dom.StatusVM{
		dom.VM_STATUS_PROVISIONING,
		dom.VM_STATUS_STARTING,
		dom.VM_STATUS_RUNNING,
		dom.VM_STATUS_STOPPING,
		dom.VM_STATUS_STOPPED,
		dom.VM_STATUS_RESTARTING,
		dom.VM_STATUS_UPDATING,
		dom.VM_STATUS_CRASHED,
		dom.VM_STATUS_ERROR,
		dom.VM_STATUS_DELETING,
		dom.VM_STATUS_NOTEXIST,
	}

	failedStatus := true
	for _, v := range availablStatuses {
		if v == vm.Status {
			failedStatus = false
			break
		}
	}
	if failedStatus {
		return fmt.Errorf("vm status isn't correct: %s", vm.Status)
	}

	return nil
}
