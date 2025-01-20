package controllers

import (
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var usermodel = new(models.UserRegistration)

type UserRegistrationController struct{}

func (u *UserRegistrationController) Register(c *gin.Context) {
	var data forms.UserModel
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please enter all Mandatory Fields", "form": data, "error": err.Error()})
		c.Abort()
		return
	}
	_, err := usermodel.Register(data)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Register successfully"})
	}
}

func (u *UserRegistrationController) OtpGeneration(c *gin.Context) {

	phonenumber := c.Param("phonenumber")

	_, err := usermodel.GenerateOtp(phonenumber)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Otp sent successfully"})
	}
}

func (u *UserRegistrationController) Login(c *gin.Context) {

	phonenumber := c.Param("phonenumber")
	otp := c.Param("otp")

	userData, err := usermodel.Login(phonenumber, otp)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return

	} else {

		c.JSON(http.StatusOK, gin.H{"message": "Otp sent successfully", "data": userData})
	}
}
