package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/vzivanovic/GOLANG_FOR_STUDENTS/db"
)

var grpcHostname string

func main() {
	flag.StringVar(&grpcHostname, "grpc-hostname", "localhost", "gRPC server hostname")
	flag.Parse()

	db.InitLocationDB()
	defer db.CloseDB()

	r := gin.Default()

	r.POST("/api/v1/location/update", func(c *gin.Context) {
		UpdateLocationHandler(c, grpcHostname, db.DB)
	})
	r.GET("/api/v1/location/search", func(c *gin.Context) {
		SearchUsersHandler(c, db.DB)
	})
	r.GET("/api/v1/location/distance", func(c *gin.Context) {
		GetDistanceHandler(c, grpcHostname, db.DB)
	})

	r.Run(":8080")
}
