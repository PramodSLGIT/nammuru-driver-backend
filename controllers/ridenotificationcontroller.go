package controllers

import (
	"log"
	"nammuru-driver-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationController struct{}

func (n *NotificationController) RideRequestController(c *gin.Context) {
	conn, err := utils.UpgradeToWebSocket(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"WebSocket connection failed": err.Error()})

		c.Abort()
		return
	}
	defer conn.Close()

	for {
		var rideRequest struct {
			DriverID string `json:"driver_id"`
			Pickup   struct {
				Latitude  float64 `json:"lat"`
				Longitude float64 `json:"lon"`
			} `json:"pickup"`
		}

		if err := conn.ReadJSON(&rideRequest); err != nil {
			log.Println("Error reading JSON:", err)
			break
		}
		
		conn.WriteJSON(gin.H{"message": "New ride request", "driver_id": rideRequest.DriverID})
	}
}
