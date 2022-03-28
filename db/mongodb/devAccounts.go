package mongodb

import (
	"devZoneDeployment/api"
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DEV_ACC_COLLECTION_NAME = "dev_accounts"

func (ds *MongoDBSource) GetDevAccounts(uid api.UserIdentity, f db.Filter) []*dom.DevAccount {
	devAccs := ds.Database.Collection(DEV_ACC_COLLECTION_NAME)
	findCursor, err := devAccs.Find(DefaulContext(), f.Compose())
	if err != nil {
		log.Fatal(err)
	}

	capacity := 1
	if f.Empty() {
		capacity = 10
	}
	res := make([]*dom.DevAccount, 0, capacity)
	for findCursor.Next(DefaulContext()) {
		devAccount := &dom.DevAccount{}
		if err := findCursor.Decode(devAccount); err != nil {
			log.Fatal(err)
		}
		res = append(res, devAccount)
	}

	return res
}

func (ds *MongoDBSource) SetDevAccounts(uid api.UserIdentity, devAccount *dom.DevAccount, isNew bool) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can change developer accounts")
	}

	devAccs := ds.Database.Collection(DEV_ACC_COLLECTION_NAME)
	if isNew {
		devAccount.Id = ds.getNewId(DEV_ACC_COLLECTION_NAME)
		res, err := devAccs.InsertOne(DefaulContext(), devAccount)
		if err != nil {
			log.Printf("Cannot insert %v, cause: %s\n", devAccount, err)
			return err
		}
		log.Println("DEV ACCOUNT INSERT (document_id):", res.InsertedID)
	} else {
		filter := new(mongoFilter)
		filter.AddEq("_id", devAccount.Id)
		res, err := devAccs.ReplaceOne(DefaulContext(), filter.Compose(), devAccount)
		if err != nil {
			log.Printf("Cannot update %v, cause: %s\n", devAccount, err)
			return err
		} else if res.ModifiedCount == 0 {
			return fmt.Errorf("the account has not been modified, because 0 accs have such id: %d", devAccount.Id)
		}
		log.Println("DEV ACCOUNT UPDATE (modified count):", res.ModifiedCount)
	}

	return nil
}

func (ds *MongoDBSource) RemoveDevAccounts(uid api.UserIdentity, devAccount *dom.DevAccount) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can delete developer accounts")
	}

	// Get app users which have a deleted Dev acc
	aus := ds.GetAppUsers(uid, ds.GetFilter("devAccountId", devAccount.Id))
	if len(aus) > 0 {
		// Cannot delete devacc which is bound with app users
		return fmt.Errorf("deleted dev account is bound with app user (name): %s", aus[0].Username)
	}

	devAccs := ds.Database.Collection(DEV_ACC_COLLECTION_NAME)
	filter := new(mongoFilter)
	filter.AddEq("_id", devAccount.Id)
	res, err := devAccs.DeleteOne(DefaulContext(), filter.Compose())
	if err != nil {
		log.Printf("Cannot remove %v, cause: %s\n", devAccount, err)
		return err
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("the account has not been deleted, because 0 accs have such id: %d", devAccount.Id)
	}

	return nil
}

func (ds *MongoDBSource) getNewId(collection string) uint {
	counters := ds.Database.Collection("counters")
	filter := (&mongoFilter{}).AddEq("_id", collection)
	setStmt := bson.D{{Key: "$inc", Value: bson.D{{Key: "val", Value: 1}}}}
	after := options.After
	upsert := true
	options := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	res := counters.FindOneAndUpdate(DefaulContext(), filter.Compose(), setStmt, &options)

	counterObj := struct {
		Id  string `bson:"_id"`
		Val uint   `bson:"val"`
	}{}
	if err := res.Decode(&counterObj); err != nil {
		log.Panic(err)
	}

	return counterObj.Val
}

// ToDo all Fields dosn't work due there aren't all fields on front
func (ds *MongoDBSource) devAccountExists(uid api.UserIdentity, devAccount *dom.DevAccount, checkAllFields bool) bool {
	filter := ds.GetFilter("_id", devAccount.Id)
	accounts := ds.GetDevAccounts(uid, filter)
	if len(accounts) != 1 {
		log.Printf("Account %#v is absent", devAccount)
		return false
	}

	if checkAllFields && !reflect.DeepEqual(devAccount, accounts[0]) {
		log.Printf("Account %#v exists, but is not equal to the instance in DB %#v", devAccount, accounts[0])
		return false
	}

	return true
}
