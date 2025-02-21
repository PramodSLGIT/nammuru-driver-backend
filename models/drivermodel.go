package models

import (
	"context"
	"fmt"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/utils"

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

func (d *DriverModel) GetNearbyDrivers(customerLat, customerLon, radius float64) ([]forms.DriverLocation, error) {
	ctx := context.TODO()

	driverIds, err := utils.RedisClient.GeoSearch(ctx, "drivers", &redis.GeoSearchQuery{
		Longitude:  customerLon,
		Latitude:   customerLat,
		Radius:     radius,
		RadiusUnit: "km",
		Sort:       "ASC",
		Count:      10,
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch drivers: %v", err)
	}

	if len(driverIds) == 0 {
		return []forms.DriverLocation{}, nil
	}

	var driverlocation []forms.DriverLocation

	for _, driverId := range driverIds {
		pos, err := utils.RedisClient.GeoPos(context.TODO(), "drivers", driverId).Result()
		if err != nil {
			fmt.Printf("Error fetching location for driver %s: %v\n", driverId, err)
			continue
		}

		if len(pos) > 0 && pos[0] != nil {
			driverlocation = append(driverlocation, forms.DriverLocation{
				ID:        driverId,
				Latitude:  pos[0].Latitude,
				Longitude: pos[0].Longitude,
			})
		}
	}
	return driverlocation, nil
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
