package mongodb

import (
	"context"
	"devZoneDeployment/db"

	"go.mongodb.org/mongo-driver/mongo"
)

var DS *MongoDBSource

// Serve MongoDB, exec CRUD
// Keep open connection to DB
// Implementation interface "DataActions"
type MongoDBSource struct {
	Database *mongo.Database
}

func UseMongoDBSource() *MongoDBSource {
	DS = newMongoDBSource()
	return DS
}

func (ds *MongoDBSource) GetEmptyFilter() db.Filter {
	return &mongoFilter{}
}

func (ds *MongoDBSource) GetFilter(field string, value interface{}) db.Filter {
	newFilter := &mongoFilter{}
	if field != "" {
		newFilter.AddEq(field, value)
	}
	return newFilter
}

func (ds *MongoDBSource) Close() {
	ds.Database.Client().Disconnect(context.TODO())
}
