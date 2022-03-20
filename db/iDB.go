package db

import (
	"devZoneDeployment/db/dom"
)

// Interface define actions which performs with database
type DataActions interface {
	// Application users section
	GetAppUsers(uid dom.UserIdentity, f Filter) []*dom.User
	GetAppUserById(uid dom.UserIdentity, userId uint) *dom.User
	GetAppUserByName(uid dom.UserIdentity, username string) *dom.User
	SetAppUser(uid dom.UserIdentity, user *dom.User, isNew bool) error
	RemoveAppUser(uid dom.UserIdentity, user *dom.User) error
	SetAppUserPassword(uid dom.UserIdentity, userId uint, pass string) error

	// Developer accounts section
	GetDevAccounts(uid dom.UserIdentity, f Filter) []*dom.DevAccount
	SetDevAccounts(uid dom.UserIdentity, acc *dom.DevAccount, isNew bool) error
	RemoveDevAccounts(uid dom.UserIdentity, acc *dom.DevAccount) error

	// Devzone actions
	GetNewOrExistsOpenVPNKey(uid dom.UserIdentity, acc *dom.DevAccount) ([]byte, error)
	ListVirtualMachines(uid dom.UserIdentity) []*dom.VM
	UpdateListVirtualMachinesFromCloud(uid dom.UserIdentity) error

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
