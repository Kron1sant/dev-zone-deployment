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

type dbParams struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Name string `yaml:"name"`
}

// Create and intialize new instance of MongoDB source (connector)
func newMongoDBSource() *MongoDBSource {
	// Get db params
	var dbParams dbParams
	if err := config.AppConfig.Sub("db").Unmarshal(&dbParams); err != nil {
		log.Fatalf("cannot get db params from config: %s", err)
	}

	// Password keeps separately
	dbParams.Pass = config.SecConfig.GetDBPassword()
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbParams.User, dbParams.Pass, dbParams.Host, dbParams.Port)
	// Create new client to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("cannot create mongoDB client: %s", err)
	} else {
		log.Printf("New MongoDB client to %s:%s\n", dbParams.Host, dbParams.Port)
	}

	// Connect to MongoDB
	if err := client.Connect(defaulContext()); err != nil {
		log.Fatalf("cannot connect to MongoDB: %s", err)
	}

	// Check connection
	if err := client.Ping(defaulContext(), nil); err != nil {
		log.Fatalf("mongoDB ping failed: %s", err)
	}

	// Prepare DataSoure
	newMongoDBSource := MongoDBSource{
		Database: client.Database(dbParams.Name),
	}

	// Check existence of App admin, otherwise create one
	newMongoDBSource.checkAdmin()

	// Return new DB connector
	return &newMongoDBSource
}

// Return the context to operation with MongoDB
func defaulContext() context.Context {
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
