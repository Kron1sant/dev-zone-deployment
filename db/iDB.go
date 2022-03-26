package db

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db/dom"
)

// Interface define actions which performs with database
type DataActions interface {
	// Application users section
	GetAppUsers(uid api.UserIdentity, f Filter) []*dom.User
	GetAppUserById(uid api.UserIdentity, userId uint) *dom.User
	GetAppUserByName(uid api.UserIdentity, username string) *dom.User
	SetAppUser(uid api.UserIdentity, user *dom.User, isNew bool) error
	RemoveAppUser(uid api.UserIdentity, user *dom.User) error
	SetAppUserPassword(uid api.UserIdentity, userId uint, pass string) error

	// Developer accounts section
	GetDevAccounts(uid api.UserIdentity, f Filter) []*dom.DevAccount
	SetDevAccounts(uid api.UserIdentity, acc *dom.DevAccount, isNew bool) error
	RemoveDevAccounts(uid api.UserIdentity, acc *dom.DevAccount) error

	// Devzone actions
	GetNewOrExistsOpenVPNKey(uid api.UserIdentity, acc *dom.DevAccount) ([]byte, error)
	ListVirtualMachines(uid api.UserIdentity) []*dom.VM
	UpdateListVirtualMachinesFromCloud(uid api.UserIdentity) error

	// Auxiliary methods
	GetEmptyFilter() Filter
	GetFilter(field string, value interface{}) Filter
	Close()
}

// Define query filter
type Filter interface {
	// Add one filtering condition
	AddEq(field string, value interface{}) Filter
	AddNEq(field string, value interface{}) Filter
	// Compose all condition at bson object containing query operators
	Compose() interface{}
	Empty() bool
}
