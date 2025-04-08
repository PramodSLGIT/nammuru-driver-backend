package controllers

import (
	"encoding/base64"
	"nammuru-driver-backend/forms"
	"nammuru-driver-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

var usermodel = new(models.UserRegistration)

type UserRegistrationController struct{}

func (u *UserRegistrationController) Register(c *gin.Context) {
	var data forms.UserRegistrationForm
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please enter all Mandatory Fields", "form": data, "error": err.Error()})
		c.Abort()
		return
	}

	var profileImage []byte
	if data.ProfileImage != "" {
		imageData, err := base64.StdEncoding.DecodeString(data.ProfileImage)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 image: " + err.Error()})
			c.Abort()
			return
		}
		if len(imageData) > 100*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image is should be 100kb"})
			c.Abort()
			return
		}
		profileImage = imageData
	}
	user := forms.UserModel{
		ProfileImage: profileImage,
		Name:         data.Name,
		Email:        data.Email,
		PhoneNumber:  data.PhoneNumber,
		Gender:       data.Gender,
		KYCData:      data.KYCData,
		Password:     data.Password, // Consider hashing it before saving!
	}
	_, err := usermodel.Register(user)
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": "Please enter all mandatory fields", "error": err.Error()})

		c.Abort()
		return

	}
	c.JSON(http.StatusOK, gin.H{"message": "Register successfully"})
}

func (u *UserRegistrationController) OtpGeneration(c *gin.Context) {

	phonenumber := c.Param("phonenumber")

	_, err := usermodel.GenerateOtp(phonenumber)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Otp sent successfully"})
}

func (u *UserRegistrationController) Login(c *gin.Context) {

	phonenumber := c.Param("phonenumber")
	otp := c.Param("otp")

	userData, err := usermodel.Login(phonenumber, otp)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

		c.Abort()
		return

	}

	c.JSON(http.StatusOK, gin.H{"message": "Otp sent successfully", "data": userData})
}

// func (u *UserRegistrationController) AddProfileImage(c *gin.Context) {
// 	var ProfileImage struct {
// 		Image       string `json:"image" bson:"image" binding:"-"`
// 		PhoneNumber string `json:"phone_number" bson:"phone_number" binding:"-"`
// 	}

// 	if err := c.BindJSON(&ProfileImage); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

// 		c.Abort()
// 		return
// 	}

// 	imageData, err := base64.StdEncoding.DecodeString(ProfileImage.Image)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Base64 image: " + err.Error()})
// 		c.Abort()
// 		return
// 	}

// 	data := forms.UserModel{
// 		ProfileImage: imageData,
// 		PhoneNumber:  ProfileImage.PhoneNumber,
// 	}

// 	_, err = usermodel.AddProfileImage(data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})

// 		c.Abort()
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Profile image added Successfully"})
// }
