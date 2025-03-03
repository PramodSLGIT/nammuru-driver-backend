package controllers

import (
	"log"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/models"
	"nammuru-driver-backend/utils"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DriverController struct{}

var drivermodel = new(models.DriverModel)

var DriverConnections = make(map[string]*websocket.Conn)
var DriverMutex = sync.Mutex{}

// WebSocket connection handler for drivers
// func (d *DriverController) DriverWebSocket(c *gin.Context) {
// 	conn, err := utils.UpgradeToWebSocket(c)
// 	if err != nil {
// 		log.Println("WebSocket connection failed:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Keep listening for messages
// 	for {
// 		messageType, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Driver disconnected:", err)
// 			utils.RemoveDriver(conn) // Remove driver from active list
// 			break                    // Exit loop
// 		}

// 		// If driver sends "exit", close connection
// 		if string(msg) == "exit" {
// 			log.Println("🚪 Driver requested to exit")
// 			utils.RemoveDriver(conn)
// 			break
// 		}

// 		log.Println("Message from driver:", string(msg))

// 		// Echo back to driver (for testing)
// 		err = conn.WriteMessage(messageType, msg)
// 		if err != nil {
// 			log.Println("Error sending message:", err)
// 			break
// 		}
// 	}
// }

func (d *DriverController) DriverLocation(c *gin.Context) {
	conn, err := utils.UpgradeToWebSocket(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"WebSocket connection failed": err.Error()})

		c.Abort()
		return
	}
	defer conn.Close()

	var location forms.DriverLocation

	err = conn.ReadJSON(&location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Location failed": err.Error()})

		c.Abort()
		return
	}

	DriverMutex.Lock()
	DriverConnections[location.ID] = conn
	DriverMutex.Unlock()

	err = drivermodel.UpdateDriverLocation(location)
	if err != nil {
		log.Println("Error updating location:", err)
	}
	defer func() {
		DriverMutex.Lock()
		delete(DriverConnections, location.ID) // Remove driver when disconnected
		DriverMutex.Unlock()
		log.Println("Driver disconnected:", location.ID)
	}()

	for {

		err := conn.ReadJSON(&location)
		if err != nil {
			log.Println("Error reading JSON:", err)
			break
		}

		err = drivermodel.UpdateDriverLocation(location)
		if err != nil {
			log.Println("Error updating location:", err)
			break
		}

		resp := map[string]string{
			"status":  "success",
			"message": "Location updated",
		}

		if err := conn.WriteJSON(resp); err != nil {
			log.Println("Error sending response:", err)
			break
		}
	}
	log.Println("Driver disconnected from WebSocket")

}

func (d *DriverController) GetNearbyDrivers(c *gin.Context) {
	var CustomerLocation struct {
		Radius    float64 `json:"radius"`
		Latitude  float64 `json:"lat"`
		Longitude float64 `json:"lon"`
	}

	if err := c.BindJSON(&CustomerLocation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "customer location is not correct"})

		c.Abort()
		return
	}
	drivers, err := drivermodel.GetNearbyDrivers(CustomerLocation.Latitude, CustomerLocation.Longitude, CustomerLocation.Radius)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get drivers"})

		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": drivers})
}

func (d *DriverController) GetAllDrivers(c *gin.Context) {
	drivers, err := drivermodel.GetAllDrivers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get drivers"})

		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": drivers})
}
