package models

import (
	"context"
	"fmt"
	"log"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Driver struct {
	Connection *websocket.Conn
	Lat        float64
	Lon        float64
}
type DriverModel struct{}

func (d *DriverModel) UpdateDriverLocation(location forms.DriverLocation) error {
	_, err := utils.RedisClient.GeoAdd(context.TODO(), "drivers", &redis.GeoLocation{
		Name:      location.ID,
		Longitude: location.Longitude,
		Latitude:  location.Latitude,
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to store driver location: %v", err)
	}

	fmt.Printf("Driver %s location updated: lat=%.6f, lon=%.6f\n", location.ID, location.Latitude, location.Longitude)
	return nil
}

func (d *DriverModel) GetNearbyDrivers(customerLocation forms.CustomerLocation) ([]forms.DriverLocation, error) {
	ctx := context.TODO()
	drivers, err := utils.RedisClient.GeoSearch(ctx, "drivers", &redis.GeoSearchQuery{
		Longitude:  customerLocation.Pickup.Longitude,
		Latitude:   customerLocation.Pickup.Latitude,
		Radius:     customerLocation.Radius,
		RadiusUnit: "km",
		Sort:       "ASC",
		Count:      10,
	}).Result()

	if err != nil || len(drivers) == 0 {
		return nil, fmt.Errorf("no available drivers")
	}

	var driverlocation []forms.DriverLocation

	for _, driverId := range drivers {
		forms.DriverMutex.Lock()
		conn, exists := forms.DriverConnections[driverId]
		forms.DriverMutex.Unlock()

		if !exists {
			log.Println("Driver", driverId, "is not online, skipping.")
			continue
		}

		err := conn.WriteJSON(gin.H{"message": "New ride request", "customer_id": customerLocation})
		if err != nil {
			log.Println("Error sending ride request to driver", driverId, ":", err)
			continue
		}

		log.Println("Ride request sent to driver", driverId)

		if WaitForDriverResponse(driverId, 10*time.Second) {
			log.Println("Driver", driverId, "accepted the ride. Stopping further requests.")
			return driverlocation, nil
		}
	}

	return nil, fmt.Errorf("no drivers accepted the ride")
}

func (d *DriverModel) GetAllDrivers() ([]forms.DriverLocation, error) {
	ctx := context.TODO()

	drivers, err := utils.RedisClient.ZRange(ctx, "drivers", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch driver IDs: %v", err)
	}

	// If no drivers are found, return empty
	if len(drivers) == 0 {
		return []forms.DriverLocation{}, nil
	}

	var driverLocations []forms.DriverLocation

	// Fetch location of each driver
	for _, driverID := range drivers {
		// Get the latitude and longitude of the driver
		pos, err := utils.RedisClient.GeoPos(ctx, "drivers", driverID).Result()
		if err != nil {
			fmt.Printf("Error fetching location for driver %s: %v\n", driverID, err)
			continue
		}

		// pos[0] contains latitude & longitude if found
		if len(pos) > 0 && pos[0] != nil {
			driverLocations = append(driverLocations, forms.DriverLocation{
				ID:        driverID,
				Latitude:  pos[0].Latitude,
				Longitude: pos[0].Longitude,
			})
		}
	}

	return driverLocations, nil
}

func WaitForDriverResponse(driverID string, timeout time.Duration) bool {
	start := time.Now()

	for time.Since(start) < timeout {
		forms.DriverMutex.Lock()
		accepted, exists := forms.RideAcceptStatus[driverID]
		if exists {
			delete(forms.RideAcceptStatus, driverID) // Remove status after reading
		}
		forms.DriverMutex.Unlock()

		if exists && accepted {
			return true
		}
		time.Sleep(1 * time.Second)
	}
	return false
}
