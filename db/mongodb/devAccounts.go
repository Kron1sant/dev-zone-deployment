package mongodb

import (
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"fmt"
	"log"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ds *MongoDBSource) GetDevAccounts(uid dom.UserIdentity, f db.Filter) []*dom.DevAccount {
	devAccs := ds.Database.Collection("dev_accounts")
	findCursor, err := devAccs.Find(defaulContext(), f.Compose())
	if err != nil {
		log.Fatal(err)
	}

	capacity := 1
	if f.Empty() {
		capacity = 10
	}
	res := make([]*dom.DevAccount, 0, capacity)
	for findCursor.Next(defaulContext()) {
		devAccount := &dom.DevAccount{}
		if err := findCursor.Decode(devAccount); err != nil {
			log.Fatal(err)
		}
		res = append(res, devAccount)
	}

	return res
}

func (ds *MongoDBSource) SetDevAccounts(uid dom.UserIdentity, devAccount *dom.DevAccount, isNew bool) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can change developer accounts")
	}

	devAccs := ds.Database.Collection("dev_accounts")
	if isNew {
		devAccount.Id = ds.getNewId("dev_accounts")
		res, err := devAccs.InsertOne(defaulContext(), devAccount)
		if err != nil {
			log.Printf("Cannot insert %v, cause: %s\n", devAccount, err)
			return err
		}
		log.Println("RESULT INSERTING:", *res)
	} else {
		filter := new(mongoFilter)
		filter.AddEq("_id", devAccount.Id)
		_, err := devAccs.ReplaceOne(defaulContext(), filter.Compose(), devAccount)
		if err != nil {
			log.Printf("Cannot update %v, cause: %s\n", devAccount, err)
			return err
		}
	}

	return nil
}

func (ds *MongoDBSource) RemoveDevAccounts(uid dom.UserIdentity, devAccount *dom.DevAccount) error {
	if !uid.IsAdmin {
		return fmt.Errorf("only admin can delete developer accounts")
	}

	// Get app users which have a deleted Dev acc
	aus := ds.GetAppUsers(uid, ds.GetFilter("devAccountId", devAccount.Id))
	if len(aus) > 0 {
		// Cannot delete devacc which is bound with app users
		return fmt.Errorf("deleted dev account is bound with app user (name): %s", aus[0].Username)
	}

	devAccs := ds.Database.Collection("dev_accounts")
	filter := new(mongoFilter)
	filter.AddEq("_id", devAccount.Id)
	_, err := devAccs.DeleteOne(defaulContext(), filter.Compose())
	if err != nil {
		log.Printf("Cannot remove %v, cause: %s\n", devAccount, err)
		return err
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
	res := counters.FindOneAndUpdate(defaulContext(), filter.Compose(), setStmt, &options)

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
func (ds *MongoDBSource) devAccountExists(uid dom.UserIdentity, devAccount *dom.DevAccount, checkAllFields bool) bool {
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
