package api

import (
	"nammuru-driver-backend/controllers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func DriverApi() {
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Authorization", "Expire", "Token"}
	config.AllowMethods = []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"}

	router.Use(cors.New(config))

	c1 := router.Group("/c1")
	{
		user := new(controllers.UserRegistrationController)
		c1.POST("/registration", user.Register)
		c1.GET("/otp/:phonenumber", user.OtpGeneration)
		c1.GET("/verify/:phonenumber/:otp", user.Login)
	}

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
	router.Run(":8020")
}
