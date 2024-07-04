package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto"
)

func UpdateLocationHandler(c *gin.Context, grpcHostname string) {
	var req LocationUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := grpc.DialContext(context.Background(), grpcHostname+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to location history microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to location history microservice"})
		return
	}
	defer conn.Close()

	client := pb.NewLocationServiceClient(conn)

	_, err = client.UpdateLocation(context.Background(), &pb.LocationUpdate{
		Username:  req.Username,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		log.Printf("Failed to update location in microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in microservice"})
		return
	}

	if err := updateLocation(req); err != nil {
		log.Printf("Failed to update location in database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "location updated"})
}

func SearchUsersHandler(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := searchUsers(req)
	c.JSON(http.StatusOK, res)
}

func GetDistanceHandler(c *gin.Context, grpcHostname string) {
	var req DistanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set EndTime to current time if not provided
	if req.EndTime.IsZero() {
		req.EndTime = time.Now()
	}

	conn, err := grpc.DialContext(context.Background(), grpcHostname+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to location history microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to location history microservice"})
		return
	}
	defer conn.Close()

	client := pb.NewLocationServiceClient(conn)

	res, err := client.GetDistance(context.Background(), &pb.DistanceRequest{
		Username:  req.Username,
		StartTime: timestamppb.New(req.StartTime),
		EndTime:   timestamppb.New(req.EndTime),
	})
	if err != nil {
		log.Printf("Failed to calculate distance in microservice: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate distance in microservice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"distance": res.Distance})
}
