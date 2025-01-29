package controllers

import (
	"encoding/base64"
	"log"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleController struct{}

var vehiclemodel = new(models.VehicleModel)

func (u *VehicleController) AddVehicleData(c *gin.Context) {
	phoneNumber := c.Param("phonenumber")
	var data forms.Vechicle

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return
	}

	_, err := vehiclemodel.AddVehicleData(data, phoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle added Successfully"})
}

func (u *VehicleController) AddVehicleImage(c *gin.Context) {
	var ProfileImage struct {
		Image    string `json:"image" bson:"image" binding:"-"`
		DriverId string `json:"driver_id" bson:"driver_id" binding:"-"`
	}

	if err := c.BindJSON(&ProfileImage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return
	}

	imageData, err := base64.StdEncoding.DecodeString(ProfileImage.Image)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 image: " + err.Error()})
		c.Abort()
		return
	}

	objId, err := primitive.ObjectIDFromHex(ProfileImage.DriverId)
	if err != nil {
		log.Println(err)
	}

	data := forms.Vechicle{
		VehicleImage: imageData,
		DriverId:     objId,
	}

	_, err = vehiclemodel.AddVehicleImage(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vehicle image added Successfully"})
}
