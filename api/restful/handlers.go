package restful

import (
	"github.com/gin-gonic/gin"
)

// Registers REST api handlers.
// Such as:
//	* GET /users - list of application users
//	* POST /users/:action - add new/edit/delete/setPassword the application user
//	* GET /accounts - list of developer accounts
//	* POST /accounts/:action - add new/edit/delete the developer account
//	* GET /wm - list of virtual machines
//	* POST /wm/:action - add update/new/edit/delete the virtual machine
//	* POST /openvpnkey - returns the OpenVPN key for the specified dev account
func AddHandlers(router gin.IRoutes) {
	// Users handlers
	router.GET("/users", GetAppUsers())
	router.POST("/users/:action", PostAppUserAction())
	// Accounts handlers
	router.GET("/accounts", GetDevAccounts())
	router.POST("/accounts/:action", PostDevAccountAction())
	// Virtual machines handlers
	router.GET("/vm", GetVMs())
	router.POST("/vm/:action", PostVMAction())
	// OpenVPN handlers
	router.POST("/openvpnkey", PostOpenVPNKey())
}
