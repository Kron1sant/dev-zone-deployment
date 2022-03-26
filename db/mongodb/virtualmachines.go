package mongodb

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/yacloudclient"
	"fmt"
	"log"
)

func (ds *MongoDBSource) ListVirtualMachines(uid api.UserIdentity) []*dom.VM {

	return ds.getVirtualMachines(uid, ds.GetEmptyFilter())
}

func (ds *MongoDBSource) UpdateListVirtualMachinesFromCloud(uid api.UserIdentity) error {
	vms := yacloudclient.ListInstances()
	for _, v := range vms {
		foundVm := ds.getVirtualMachines(uid, ds.GetFilter("_id", v.Id))
		var vm *dom.VM
		var isNew bool
		if len(foundVm) == 0 {
			vm = &dom.VM{Id: v.Id}
			isNew = true
		} else {
			vm = foundVm[0]
			isNew = false
		}

		vm.Name = v.Name
		vm.Description = v.Description
		vm.Params = "params_foobaz"
		vm.Status = yacloudclient.InstanseStatus(v.Status)
		ds.SetVirtualMachine(uid, vm, isNew)
	}

	return nil
}

func (ds *MongoDBSource) getVirtualMachines(uid api.UserIdentity, f db.Filter) []*dom.VM {
	virtMachines := ds.Database.Collection("virtual_machines")
	findCursor, err := virtMachines.Find(DefaulContext(), f.Compose())
	if err != nil {
		log.Fatal(err)
	}

	capacity := 1
	if f.Empty() {
		capacity = 10
	}
	res := make([]*dom.VM, 0, capacity)
	for findCursor.Next(DefaulContext()) {
		vm := &dom.VM{}
		if err := findCursor.Decode(vm); err != nil {
			log.Fatal(err)
		}
		res = append(res, vm)
	}

	return res
}

func (ds *MongoDBSource) SetVirtualMachine(uid api.UserIdentity, vm *dom.VM, isNew bool) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can change virtual machine")
	}

	virtMachines := ds.Database.Collection("virtual_machines")
	if isNew {
		res, err := virtMachines.InsertOne(DefaulContext(), vm)
		if err != nil {
			log.Printf("Cannot insert %v, cause: %s\n", vm, err)
			return err
		}
		log.Println("RESULT INSERTING:", *res)
	} else {
		filter := ds.GetFilter("_id", vm.Id)
		_, err := virtMachines.ReplaceOne(DefaulContext(), filter.Compose(), vm)
		if err != nil {
			log.Printf("Cannot update %v, cause: %s\n", vm, err)
			return err
		}
	}

	return nil
}
