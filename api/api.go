package api

import (
	"log"
	"nammuru-driver-backend/controllers"
	"nammuru-driver-backend/jwtauth"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
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

	authMiddleware := jwtauth.InitJwt()
	router.POST("d1/login", authMiddleware.LoginHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	p1 := router.Group("/d1")
	{
		p1.POST("/logout", authMiddleware.LogoutHandler)
		var user = new(controllers.UserRegistrationController)
		p1.POST("/register", user.Register)
	}

	d1 := router.Group("/d1")
	d1.GET("/refresh_token", authMiddleware.RefreshHandler)
	d1.Use(authMiddleware.MiddlewareFunc())
	{
		//user
		user := new(controllers.UserRegistrationController)
		d1.GET("/otp/:phonenumber", user.OtpGeneration)
		// d1.POST("/profileimage", user.AddProfileImage)

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
