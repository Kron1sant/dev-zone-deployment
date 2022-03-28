package mongodb

import (
	"context"
	"devZoneDeployment/config"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var timeout int = 10 // mongodb connection timeout. Default 10 seconds

// Create and intialize new instance of MongoDB source (connector)
func newMongoDBSource() *MongoDBSource {
	// Composing a connection string
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", config.GetDBUser(), config.GetDBPass(), config.GetDBHost(), config.GetDBPort())
	// Creating new client to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("cannot create mongoDB client: %s", err)
	} else {
		log.Printf("New MongoDB client to %s:%s\n", config.GetDBHost(), config.GetDBPort())
	}

	// Connect to MongoDB
	if err := client.Connect(DefaulContext()); err != nil {
		log.Fatalf("cannot connect to MongoDB: %s", err)
	}

	// Check connection
	if err := client.Ping(DefaulContext(), nil); err != nil {
		log.Fatalf("mongoDB ping failed: %s", err)
	}

	// Prepare DataSoure
	newMongoDBSource := MongoDBSource{
		Database: client.Database(config.GetDBName()),
	}

	// Check existence of App admin, otherwise create one
	newMongoDBSource.checkAdmin()

	// Return new DB connector
	return &newMongoDBSource
}

// Return the context to operation with MongoDB
func DefaulContext() context.Context {
	d := time.Duration(timeout) * time.Second
	ctx, close := context.WithTimeout(context.Background(), d)
	// Check context will close after d seconds
	go func() {
		select {
		case <-ctx.Done():
			// Do nothing - context close in another place
		case <-time.After(d):
			close()
		}
	}()
	return ctx
}
