package restful

import (
	"github.com/gin-gonic/gin"
)

// Declare all functional handlers.
// Such as:
//	* GET /users - list of application users
//	* POST /users/:action - add new/edit/delete/setPassword application user
//	* GET /accounts - list of developer accounts
//	* POST /accounts/:action - add new/edit/delete developer account
//	* POST /openvpnkey - returns the OpenVPN key for the specified dev account
func AddHandlers(router gin.IRoutes) {
	router.GET("/users", GetAppUsers())
	// action can be add, edit, delete or setPassword
	router.POST("/users/:action", PostAppUserAction())

	router.GET("/accounts", GetDevAccounts())
	// action can be add, edit or delete
	router.POST("/accounts/:action", PostDevAccountAction())

	router.GET("/vm", ListVM())
	// action can be update
	router.POST("/vm/:action", PostVMAction())
	router.POST("/openvpnkey", PostOpenVPNKey())
}
