package main

import (
	"context"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/vzivanovic/GOLANG_FOR_STUDENTS/db"
	pb "github.com/vzivanovic/GOLANG_FOR_STUDENTS/proto"
)

type LocationUpdateRequest struct {
	Username  string  `json:"username" binding:"required,min=4,max=16,alphanum"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

type SearchRequest struct {
	Latitude  float64 `form:"latitude" binding:"required"`
	Longitude float64 `form:"longitude" binding:"required"`
	Radius    float64 `form:"radius" binding:"required"`
	Page      int     `form:"page" binding:"required,default=1"`
	Size      int     `form:"size" binding:"required,default=10"`
}

type DistanceRequest struct {
	Username  string    `form:"username" binding:"required,min=4,max=16,alphanum"`
	StartTime time.Time `form:"start_time" binding:"required"`
	EndTime   time.Time `form:"end_time" binding:"required"`
}

type UserLocation struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SearchResponse struct {
	Users      []UserLocation `json:"users"`
	Page       int            `json:"page"`
	TotalPages int            `json:"total_pages"`
	TotalUsers int            `json:"total_users"`
}

func updateLocation(req LocationUpdateRequest) error {
	_, err := db.DB.Exec("INSERT OR REPLACE INTO user_locations (username, latitude, longitude) VALUES (?, ?, ?)",
		req.Username, req.Latitude, req.Longitude)
	return err
}

func searchUsers(req SearchRequest) SearchResponse {
	rows, err := db.DB.Query("SELECT username, latitude, longitude FROM user_locations")
	if err != nil {
		log.Fatalf("Failed to query database: %v", err)
	}
	defer rows.Close()

	var users []UserLocation
	for rows.Next() {
		var user UserLocation
		if err := rows.Scan(&user.Username, &user.Latitude, &user.Longitude); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		if distance(req.Latitude, req.Longitude, user.Latitude, user.Longitude) <= req.Radius {
			users = append(users, user)
		}
	}

	totalUsers := len(users)
	totalPages := int(math.Ceil(float64(totalUsers) / float64(req.Size)))
	start := (req.Page - 1) * req.Size
	end := start + req.Size
	if end > totalUsers {
		end = totalUsers
	}
	paginatedUsers := users[start:end]

	return SearchResponse{
		Users:      paginatedUsers,
		Page:       req.Page,
		TotalPages: totalPages,
		TotalUsers: totalUsers,
	}
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in kilometers
	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func main() {
	db.InitLocationDB()

	r := gin.Default()

	r.POST("/api/v1/location/update", func(c *gin.Context) {
		var req LocationUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Connect to location history microservice
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to location history microservice: %v", err)
		}
		defer conn.Close()

		client := pb.NewLocationServiceClient(conn)

		_, err = client.UpdateLocation(context.Background(), &pb.LocationUpdate{
			Username:  req.Username,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
			return
		}

		// Update database
		if err := updateLocation(req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location in database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "location updated"})
	})

	r.GET("/api/v1/location/search", func(c *gin.Context) {
		var req SearchRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res := searchUsers(req)
		c.JSON(http.StatusOK, res)
	})

	r.GET("/api/v1/location/distance", func(c *gin.Context) {
		var req DistanceRequest
		if err := c.ShouldBindQuery(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Connect to location history microservice
		conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to location history microservice: %v", err)
		}
		defer conn.Close()

		client := pb.NewLocationServiceClient(conn)

		res, err := client.GetDistance(context.Background(), &pb.DistanceRequest{
			Username:  req.Username,
			StartTime: timestamppb.New(req.StartTime),
			EndTime:   timestamppb.New(req.EndTime),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate distance"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"distance": res.Distance})
	})

	r.Run(":8080")
}
