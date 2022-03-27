package mongodb

import (
	"devZoneDeployment/api"
	"devZoneDeployment/config"
	"devZoneDeployment/db"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/utils"
	"fmt"
	"log"
)

func (ds *MongoDBSource) GetAppUsers(uid api.UserIdentity, f db.Filter) []*dom.User {
	// Add a filter which control access
	// Non admin user can see only own record
	if !uid.IsAdmin {
		f.AddEq("_id", uid.Id)
	}

	return ds.getAppUsersFiltered(f)
}

func (ds *MongoDBSource) SetAppUser(uid api.UserIdentity, appUser *dom.User, isNew bool) error {
	if err := ds.isAppUserEditingAvailable(uid, appUser, isNew); err != nil {
		return err
	}

	appUsers := ds.Database.Collection("app_users")
	if isNew {
		appUser.Id = ds.getNewId("app_users")
		res, err := appUsers.InsertOne(DefaulContext(), appUser)
		if err != nil {
			log.Printf("Cannot insert %v, cause: %s\n", appUser, err)
			return err
		}
		log.Println("RESULT INSERTING:", *res)
	} else {
		filter := new(mongoFilter)
		filter.AddEq("_id", appUser.Id)
		res, err := appUsers.ReplaceOne(DefaulContext(), filter.Compose(), appUser)
		if err != nil {
			log.Printf("Cannot update %v, cause: %s\n", appUser, err)
			return err
		} else if res.ModifiedCount == 0 {
			return fmt.Errorf("the account has not been modified, because 0 accs have such id: %d", appUser.Id)
		}
	}

	return nil
}

func (ds *MongoDBSource) RemoveAppUser(uid api.UserIdentity, appUser *dom.User) error {
	if !uid.IsAdmin && uid.Id != appUser.Id {
		return fmt.Errorf("only admin can delete another user")
	}

	appUsers := ds.Database.Collection("app_users")
	filter := new(mongoFilter)
	filter.AddEq("_id", appUser.Id)
	res, err := appUsers.DeleteOne(DefaulContext(), filter.Compose())
	if err != nil {
		log.Printf("Cannot remove %v, cause: %s\n", appUser, err)
		return err
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("the user has not been deleted, because 0 users have such id: %d", appUser.Id)
	}

	return nil
}

func (ds *MongoDBSource) SetAppUserPassword(uid api.UserIdentity, userId uint, password string) error {
	if !uid.IsAdmin && uid.Id != userId {
		return fmt.Errorf("only admin can change another user passwords")
	}

	// Find app user by id
	user := ds.GetAppUserById(uid, userId)
	if user == nil {
		return fmt.Errorf("cannot find user by id %d", userId)
	}
	user.Password = utils.HashAndSaltPassword(password)
	return ds.SetAppUser(uid, user, false)
}

func (ds *MongoDBSource) GetAppUserById(uid api.UserIdentity, userId uint) *dom.User {
	filter := ds.GetFilter("_id", userId)
	users := ds.getAppUsersFiltered(filter)
	if len(users) == 1 {
		return users[0]
	} else {
		log.Printf("User didn't find by id %d", userId)
		return nil
	}
}

func (ds *MongoDBSource) GetAppUserByName(uid api.UserIdentity, username string) *dom.User {
	filter := ds.GetFilter("username", username)
	users := ds.getAppUsersFiltered(filter)
	if len(users) == 1 {
		return users[0]
	} else {
		return nil
	}
}

func (ds *MongoDBSource) getAppUsersFiltered(f db.Filter) []*dom.User {
	appUsers := ds.Database.Collection("app_users")
	findCursor, err := appUsers.Find(DefaulContext(), f.Compose())
	if err != nil {
		log.Fatal(err)
	}

	capacity := 10
	if !f.Empty() {
		capacity = 1
	}
	res := make([]*dom.User, 0, capacity)
	for findCursor.Next(DefaulContext()) {
		appUser := &dom.User{}
		if err := findCursor.Decode(appUser); err != nil {
			log.Fatal(err)
		}
		res = append(res, appUser)
	}

	return res
}

func (ds *MongoDBSource) checkAdmin() {
	// Get default admin from config
	// Default admin is needed to provide app authentication (e.g. RESTful api)
	adminFromConfig := config.GetDefaultAdmin()
	defaultAdmin := dom.User{
		Username: adminFromConfig.Username,
		EMail:    adminFromConfig.Email,
		Password: utils.HashAndSaltPassword(adminFromConfig.Password),
		IsAdmin:  true,
	}

	// Check admin presence, otherwise add one
	uid := api.IdentityFromUser(&defaultAdmin)
	if ds.GetAppUserByName(uid, defaultAdmin.Username) == nil {
		if err := ds.SetAppUser(uid, &defaultAdmin, true); err != nil {
			log.Fatalf("cannot create app admin in MongoDB: %s", err)
		} else {
			log.Printf("Add new app admin %q\n", defaultAdmin.Username)
		}
	}
}

func (ds *MongoDBSource) isAppUserEditingAvailable(uid api.UserIdentity, appUser *dom.User, isNew bool) error {
	isAdminSession := uid.IsAdmin
	userIdOfSession := uid.Id
	var existingUser *dom.User = nil
	if !isNew {
		existingUser = ds.GetAppUserById(uid, appUser.Id)
	}

	// Rigths checking
	if !isAdminSession {
		if isNew {
			// Non admin user can't add new users
			return fmt.Errorf("only admin can add a new user")
		} else if appUser.IsAdmin {
			// Non admin user can't set admin rights himself
			return fmt.Errorf("only admin can set admin rights")
		} else if userIdOfSession != appUser.Id {
			// Non admin is not allowed to change another user
			return fmt.Errorf("current user don't have access to app user id %d ", appUser.Id)
		} else if existingUser.Username != appUser.Username {
			// Non admin are not allowed to change own username
			return fmt.Errorf("changing username is not allowed: %s -> %s", existingUser.Username, appUser.Username)
		} else if appUser.HasDevAccount && existingUser.DevAccountId != appUser.DevAccountId {
			// Non admin are not allowed to change bound devaccount
			return fmt.Errorf("changing a bound devaccont is not allowed: %d -> %d", existingUser.DevAccountId, appUser.DevAccountId)
		}
	}

	// Consistency and validation checking
	userByName := ds.GetAppUserByName(uid, appUser.Username)
	if userByName != nil && userByName.Username == appUser.Username && userByName.Id != appUser.Id {
		// There is another user with the same name in DB -> his is unacceptable
		return fmt.Errorf("username %s is not uniq", userByName.Username)
	}

	if !isNew && existingUser == nil {
		// The user is being updated is absent
		return fmt.Errorf("the user being updated with id %d does not exist", userByName.Id)
	}

	if appUser.HasDevAccount {
		devAccs := ds.GetDevAccounts(uid, ds.GetFilter("_id", appUser.DevAccountId))
		if len(devAccs) == 0 {
			// Bad reference to dev account
			return fmt.Errorf("bad reference to dev account (by id): %d", appUser.DevAccountId)
		} else {
			// Get app users which have the same Dev acc
			aus := ds.GetAppUsers(uid, ds.GetFilter("devAccountId", appUser.DevAccountId).AddNEq("_id", appUser.Id))
			if len(aus) > 0 {
				// Specified devaccount has already bound to application user
				return fmt.Errorf("specified developer account has already bound to an application user (dev account id): %d", appUser.DevAccountId)
			}

		}
	}

	return nil
}
