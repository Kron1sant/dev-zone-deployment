package main

import (
	"context"
	"devZoneDeployment/api/restful"
	"devZoneDeployment/config"
	"devZoneDeployment/db/dom"
	"devZoneDeployment/db/mongodb"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Create DB connector and set global var
	dataSource := mongodb.UseMongoDBSource()
	defer dataSource.Close()

	// Create Web-server by Gin
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsHandler())
	authSite := restful.AddAuthMiddleware(router)
	restful.AddHandlers(authSite)

	port := config.AppConfig.GetString("service.port")
	//log.Fatal(http.ListenAndServeTLS(":"+port, "certs/localhost.crt", "certs/localhost.key", router))
	switch gin.Mode() {
	case "debug":
		router.Run(":" + port)
	case "release":
		log.Fatal(autotls.RunWithManager(router, &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(""), // ToDo
			Cache:      autocert.DirCache("certs/.cache"),
		}))
	case "test":
		log.Fatal("There is not any tests, yet:(")
	default:
		log.Fatalf("unknown gin mode %s", gin.Mode())
	}
}

// Handle CORS headers
func corsHandler() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func test(ds *mongodb.MongoDBSource) {
	appUsers := ds.Database.Collection("app_users")

	findCursor, err := appUsers.Find(context.TODO(), bson.D{bson.E{"_id", bson.M{"$ne": 8}}})
	if err != nil {
		log.Fatal(err)
	}

	res := make([]*dom.User, 0, 1)
	for findCursor.Next(context.TODO()) {
		appUser := &dom.User{}
		if err := findCursor.Decode(appUser); err != nil {
			log.Fatal(err)
		}
		res = append(res, appUser)
	}

	fmt.Println(res)
	log.Fatal("asdas")
}
