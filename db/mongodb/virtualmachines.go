package mongodb

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/yacloudclient"
	"fmt"
	"log"
)

const VM_COLLECTION_NAME = "virtual_machines"

func (ds *MongoDBSource) ListVirtualMachines(uid api.UserIdentity) []*dom.VM {

	return ds.getVirtualMachines(uid, ds.GetEmptyFilter())
}

func (ds *MongoDBSource) SetVirtualMachine(uid api.UserIdentity, vm *dom.VM, isNew bool) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can change virtual machine")
	}

	vmColls := ds.Database.Collection(VM_COLLECTION_NAME)
	if isNew {
		res, err := vmColls.InsertOne(DefaulContext(), vm)
		if err != nil {
			log.Printf("Cannot insert %v, cause: %s\n", vm, err)
			return err
		}
		log.Println("VM INSERT (document_id):", res.InsertedID)
	} else {
		filter := ds.GetFilter("_id", vm.Id)
		res, err := vmColls.ReplaceOne(DefaulContext(), filter.Compose(), vm)
		if err != nil {
			log.Printf("Cannot update %v, cause: %s\n", vm, err)
			return err
		} else if res.ModifiedCount == 0 {
			return fmt.Errorf("the vm has not been modified, because 0 vms have such id: %s", vm.Id)
		}
		log.Println("VM UPDATE (modified count)", res.ModifiedCount)
	}

	return nil
}

func (ds *MongoDBSource) RemoveVirtualMachine(uid api.UserIdentity, vm *dom.VM) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can delete virtual machine")
	}

	// Get dev accounts that have the deleted VM
	accs := ds.GetDevAccounts(uid, ds.GetFilter("_id", vm.DevAccountId))
	if len(accs) > 0 {
		// Cannot delete VM which is bound with devacc
		return fmt.Errorf("deleted virtual machine is bound with dev account (name): %s", accs[0].Username)
	}

	vmColls := ds.Database.Collection(VM_COLLECTION_NAME)
	filter := new(mongoFilter)
	filter.AddEq("_id", vm.Id)
	res, err := vmColls.DeleteOne(DefaulContext(), filter.Compose())
	if err != nil {
		log.Printf("Cannot remove %v, cause: %s\n", vm, err)
		return err
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("the vm has not been deleted, because 0 vms have such id: %s", vm.Id)
	}

	return nil
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
		vm.Status = dom.StatusVM(yacloudclient.InstanseStatus(v.Status))
		ds.SetVirtualMachine(uid, vm, isNew)
	}

	return nil
}

func (ds *MongoDBSource) getVirtualMachines(uid api.UserIdentity, f db.Filter) []*dom.VM {
	virtMachines := ds.Database.Collection(VM_COLLECTION_NAME)
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
