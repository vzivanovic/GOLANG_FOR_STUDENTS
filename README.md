# Location Management System

## Introduction

This project is a Location Management System composed of two microservices: `location-history` and `location-management`. The system provides functionalities to update user locations, search for users within a radius, and calculate the distance traveled by a user within a specific time frame.


## Prerequisites

- Go 1.18 or later
- Protocol Buffers (protoc) version 3.0.0 or later
- Git
- SQLite3

## Setup

### 1. Clone the repository

```sh
git clone https://github.com/yourusername/GOLANG_FOR_STUDENTS.git
cd GOLANG_FOR_STUDENTS
```
### 2. Install the dependencies
```sh
cd db
go mod tidy
cd ..

cd location-history
go mod tidy
cd ..

cd location-management
go mod tidy
cd ..

protoc --go_out=. --go-grpc_out=. location.proto
```
## Run the program

It requires to run two terminals, one to run location-history and the other to run location-management.

```sh
cd location-history
go run main.go
The service will start on port: '50051'.

cd location-management
go run main.go
The service will start on port: '8080'.
```
## API Endpoints
# 1. Update location
    - URL: '/api/v1/location/update'
    - Method: 'POST'
    - Request body:
        {
            "username":"testuser",
            "latitude": 37.7749,
            "longitude": -122.4194
        }
    - Response:
        {
            "status": "location updated"
        }

# 2. Search users
    - URL: '/api/v1/location/search'
    - Method: 'GET'
    - Query parameters:
        - latitude: Latitude of the center point.
        - longitude: Longitude of the center point.
        - radius: Search radius in kilometers.
        - page: Page number (default is 1).
        - size: Number of results per page (default is 10).
    - Response :
        {
            "users": [
                {
                    "username": "testuser",
                    "latitude": 37.7749,
                    "longitude": -122.4194
                }
            ],
            "page": 1,
            "total_pages": 1,
            "total_users": 1
        }
# 3. Get distance
    - URL: 'api/v1/location/distance'
    - Method: 'GET'
    - Query parameters:
        - 'username': Username of the user
        - 'start_time': Start time in ISO 8601 format
        - 'end_time': End time in ISO 8601 format
    - Response:
        {
            "distance": 12.34
        }
