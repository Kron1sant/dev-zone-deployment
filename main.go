package main

import (
	"devZoneDeployment/api/restful"
	"devZoneDeployment/config"
	"devZoneDeployment/db/mongodb"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize an app configuration by file and env vars
	config.Init()

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

	switch gin.Mode() {
	case "debug":
		router.Run(":" + config.GetAppPort())
	case "release":
		// ToDo
	case "test":
		// ToDo
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
