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

	d1 := router.Group("/d1")
	{
		//user
		user := new(controllers.UserRegistrationController)
		d1.POST("/registration", user.Register)
		d1.GET("/otp/:phonenumber", user.OtpGeneration)
		d1.GET("/verify/:phonenumber/:otp", user.Login)
		d1.POST("/profileimage", user.AddProfileImage)

		//vehicle
		vehicle := new(controllers.VehicleController)
		d1.POST("/vehicledata/:phonenubmer", vehicle.AddVehicleData)
		d1.POST("/vehicleimage", vehicle.AddVehicleImage)

		//connect to websocket
		driver := new(controllers.DriverController)
		d1.GET("/ws/driverlocation", driver.DriverLocation)
		d1.GET("/getnearbydrivers", driver.GetNearbyDrivers)
		d1.GET("/getalldrivers", driver.GetAllDrivers)
	}

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
	router.Run(":8020")
}
