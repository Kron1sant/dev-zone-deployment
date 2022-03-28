package sessions

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/mongodb"
	"devZoneDeployment/db/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type contextSession struct {
	ctx *gin.Context
	ds  db.DataActions
	ui  api.UserIdentity
}

func New(ctx *gin.Context) *contextSession {
	ds := mongodb.DS
	ui := api.IdentityFromContext(ctx)

	return &contextSession{
		ctx: ctx,
		ds:  ds,
		ui:  ui,
	}
}

func NewGuest(ctx *gin.Context) *contextSession {
	ds := mongodb.DS
	ui := api.IdentityFromContext(ctx)
	ui.Guest = true

	return &contextSession{
		ctx: ctx,
		ds:  ds,
		ui:  ui,
	}
}

func (sctx *contextSession) GetDevAccounts() []*dom.DevAccount {
	userid := sctx.ctx.Query("accountid")
	filter := sctx.ds.GetEmptyFilter()
	if userid != "" {
		val, err := strconv.Atoi(userid)
		if err != nil {
			return make([]*dom.DevAccount, 0)
		}
		filter = sctx.ds.GetFilter("_id", val)
	}
	return sctx.ds.GetDevAccounts(sctx.ui, filter)
}

func (sctx *contextSession) AddDevAccount(acc *dom.DevAccount) error {
	return sctx.ds.SetDevAccounts(sctx.ui, acc, true)
}

func (sctx *contextSession) EditDevAccount(acc *dom.DevAccount) error {
	return sctx.ds.SetDevAccounts(sctx.ui, acc, false)
}

func (sctx *contextSession) DeleteDevAccount(acc *dom.DevAccount) error {
	return sctx.ds.RemoveDevAccounts(sctx.ui, acc)
}

func (sctx *contextSession) GetAppUsers() []*dom.User {
	userid := sctx.ctx.Query("userid")
	filter := sctx.ds.GetEmptyFilter()
	if userid != "" {
		val, err := strconv.Atoi(userid)
		if err != nil {
			return make([]*dom.User, 0)
		}
		filter = sctx.ds.GetFilter("_id", val)
	}
	return sctx.ds.GetAppUsers(sctx.ui, filter)
}

func (sctx *contextSession) AddAppUser(user *dom.User) error {
	return sctx.ds.SetAppUser(sctx.ui, user, true)
}

func (sctx *contextSession) EditAppUser(user *dom.User) error {
	// While user editing password isn't transmitted
	// We need to extract an actual password and set it to the current user
	// ToDo: other decision:
	// 	don't write password while updating DB if there is not "SetAppUserPassword" operation
	actualUser := sctx.ds.GetAppUserById(sctx.ui, user.Id)
	if actualUser == nil {
		return fmt.Errorf("user is absent in the base")
	}
	user.Password = actualUser.Password

	return sctx.ds.SetAppUser(sctx.ui, user, false)
}

func (sctx *contextSession) DeleteAppUser(user *dom.User) error {
	return sctx.ds.RemoveAppUser(sctx.ui, user)
}

func (sctx *contextSession) SetAppUserPassword(userId uint, pass string) error {
	return sctx.ds.SetAppUserPassword(sctx.ui, userId, pass)
}

func (sctx *contextSession) Authenticate(username string, password string) (*dom.User, error) {
	user := sctx.ds.GetAppUserByName(sctx.ui, username)
	if user == nil {
		return nil, fmt.Errorf("username is absent in the base")
	}

	if utils.CheckPassword(password, user.Password) {
		// success
		return user, nil
	} else {
		return nil, fmt.Errorf("wrong password")
	}
}

func (sctx *contextSession) GetNewOrExistsOpenVPNKey(acc *dom.DevAccount) ([]byte, error) {
	return sctx.ds.GetNewOrExistsOpenVPNKey(sctx.ui, acc)
}

func (sctx *contextSession) ListVirtualMachines() []*dom.VM {
	return sctx.ds.ListVirtualMachines(sctx.ui)
}

func (sctx *contextSession) AddVirtualMachine(vm *dom.VM) error {
	return sctx.ds.SetVirtualMachine(sctx.ui, vm, true)
}

func (sctx *contextSession) EditVirtualMachine(vm *dom.VM) error {
	return sctx.ds.SetVirtualMachine(sctx.ui, vm, false)
}

func (sctx *contextSession) DeleteVirtualMachine(vm *dom.VM) error {
	return sctx.ds.RemoveVirtualMachine(sctx.ui, vm)
}

func (sctx *contextSession) UpdateListVirtualMachinesFromCloud() error {
	return sctx.ds.UpdateListVirtualMachinesFromCloud(sctx.ui)
}
