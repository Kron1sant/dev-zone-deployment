package restful

import (
	"bytes"
	"devZoneDeployment/db/dom"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TEST_VIRTUAL_MACHINE = dom.VM{
	Id:            "aa24234bde230098",
	Name:          "TestVM",
	HasDevAccount: false,
	Description:   "FooBax VM",
	Status:        dom.VM_STATUS_RUNNING,
}

func TestListVM(t *testing.T) {

}

func TestPostVMAction(t *testing.T) {
	assert := assert.New(t)
	// The order of the next steps is important
	// Test adding a failed virtual machine
	stepAddFailedVM(assert)
	// Test adding a virtual machine
	testVM := stepAddCorrectVM(assert)
	// Test editing an absent virtual machine
	stepEditFailedVM(assert, testVM)
	// Test editing a virtual machine
	stepEditCorrectVM(assert, testVM)
	// Test deleting an absent virtual machine
	stepDeleteFailedVM(assert, testVM)
	// Test deleting a virtual machine
	stepDeleteCorrectVM(assert, testVM)
	// Test executing a bad action
	stepBadVMAction(assert, testVM)
}

func stepAddFailedVM(assert *assert.Assertions) {
	wrongVM := dom.VM{
		Name:   "Wrong virt machine",
		Status: "bad status",
	}
	rr := proccessVMAction("add", wrongVM)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepAddCorrectVM(assert *assert.Assertions) dom.VM {
	rr := proccessVMAction("add", TEST_VIRTUAL_MACHINE)
	// Read response a body
	createdVM := dom.VM{}
	if err := json.Unmarshal(rr.Body.Bytes(), &createdVM); err != nil {
		assert.FailNow("The response body must contain the data of the virtual machine", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(TEST_VIRTUAL_MACHINE.Name, createdVM.Name, "Created VM's name must be equal to the source")
	assert.Equal(TEST_VIRTUAL_MACHINE.Status, createdVM.Status, "Created VM's status must be equal to the source")
	assert.Equal(TEST_VIRTUAL_MACHINE.Description, createdVM.Description, "Created VM's Description must be equal to the source")
	assert.Equal(TEST_VIRTUAL_MACHINE.HasDevAccount, createdVM.HasDevAccount, "Created VM's HasDevAccount must be equal to the source")
	assert.Equal(TEST_VIRTUAL_MACHINE.DevAccountId, createdVM.DevAccountId, "Created VM's DevAccountId must be equal to the source")

	return createdVM
}

func stepEditFailedVM(assert *assert.Assertions, testVM dom.VM) {
	wrongVM := testVM
	wrongVM.Id = "00000000000000"
	rr := proccessVMAction("edit", wrongVM)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepEditCorrectVM(assert *assert.Assertions, testVM dom.VM) {
	// Change e-mail
	testVM.Description = "New description"
	rr := proccessVMAction("edit", testVM)
	// Read response a body
	editedVM := dom.VM{}
	if err := json.Unmarshal(rr.Body.Bytes(), &editedVM); err != nil {
		assert.FailNow("The response body must contain the data of the edited VM", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testVM.Id, editedVM.Id, "Edited VM's Id must be equal to the source")
	assert.Equal(testVM.Name, editedVM.Name, "Edited VM's name must be equal to the source")
	assert.Equal(testVM.Status, editedVM.Status, "Edited VM's status must be equal to the source")
	assert.Equal(testVM.Description, editedVM.Description, "Edited VM's Description must be equal to the source")
	assert.Equal(testVM.HasDevAccount, editedVM.HasDevAccount, "Edited VM's HasDevAccount must be equal to the source")
	assert.Equal(testVM.DevAccountId, editedVM.DevAccountId, "Edited VM's DevAccountId must be equal to the source")
}

func stepDeleteFailedVM(assert *assert.Assertions, testVM dom.VM) {
	wrongVM := testVM
	wrongVM.Id = "0000000000000"
	rr := proccessVMAction("delete", wrongVM)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func stepDeleteCorrectVM(assert *assert.Assertions, testVM dom.VM) {
	rr := proccessVMAction("delete", testVM)
	// Read response a body
	deletedVM := dom.VM{}
	if err := json.Unmarshal(rr.Body.Bytes(), &deletedVM); err != nil {
		assert.FailNow("The response body must contain the data of the deleted user", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusCreated, rr.Code, "The response code must be 201")
	assert.Equal(testVM.Id, deletedVM.Id, "Deleted VM's Id must be equal to the source")
	assert.Equal(testVM.Name, deletedVM.Name, "Deleted VM's name must be equal to the source")
}

func stepBadVMAction(assert *assert.Assertions, testVM dom.VM) {
	rr := proccessVMAction("bad_action", testVM)
	// Read response a body
	errorResponce := struct{ Error string }{}
	if err := json.Unmarshal(rr.Body.Bytes(), &errorResponce); err != nil {
		assert.FailNow("The response body must contain the error data", err)
	}
	// Check the code and some fields of the data
	assert.Equal(http.StatusBadRequest, rr.Code, "The response code must be 400")
	assert.NotEmpty(errorResponce.Error, "Error must not be empty")
}

func proccessVMAction(action string, acc dom.VM) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	ctx := prepareGinContext(w)
	setAction(ctx, action)
	b, _ := json.Marshal(acc)
	ctx.Request = httptest.NewRequest("POST", "/vm/"+action, bytes.NewReader(b))
	PostVMAction()(ctx)
	return w
}
