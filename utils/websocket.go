package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader (to upgrade HTTP to WebSocket)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Store connected drivers
// var drivers = make(map[*websocket.Conn]*models.Driver)
// var mutex = sync.Mutex{}

// Upgrade HTTP connection to WebSocket
func UpgradeToWebSocket(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return nil, err
	}
	return conn, nil
}

// Add driver connection
// func AddDriver(conn *websocket.Conn, lat, lon float64) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	drivers[conn] = &models.Driver{
// 		Connection: conn,
// 		Lat:        lat,
// 		Lon:        lon,
// 	}
// 	fmt.Println("Driver connected")
// }

// // Remove driver connection
// func RemoveDriver(conn *websocket.Conn) {
// 	mutex.Lock()
// 	defer mutex.Unlock()
// 	// Check if driver exists before deleting
// 	if _, exists := drivers[conn]; exists {
// 		delete(drivers, conn)
// 		fmt.Println("Driver disconnected and removed")
// 		conn.Close() // Ensure connection is closed
// 	}
// }

// func AssignNearestDriver(rideRequest string, customerLat, customerLon float64) {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	var nearestDriverConn *websocket.Conn
// 	var minDistance float64 = 99999999 // Initialize with a large value

// 	// Iterate through all drivers
// 	drivers.Range(func(key, value interface{}) bool {
// 		driver := value.(*models.Driver)

// 		// Calculate the distance between customer and driver
// 		distance := HaversineDistance(customerLat, customerLon, driver.Lat, driver.Lon)

// 		// Assign the driver if they are the nearest so far
// 		if distance < minDistance {
// 			minDistance = distance
// 			nearestDriverConn = driver.Connection
// 		}
// 		return true
// 	})

// 	// If a driver is found, send the ride request only to them
// 	if nearestDriverConn != nil {
// 		err := nearestDriverConn.WriteMessage(websocket.TextMessage, []byte(rideRequest))
// 		if err != nil {
// 			fmt.Println("❌ Error sending to driver:", err)
// 			nearestDriverConn.Close()
// 			drivers.Delete(nearestDriverConn) // Remove disconnected driver
// 		} else {
// 			fmt.Println("🚖 Ride request sent to nearest driver!")
// 		}
// 	} else {
// 		fmt.Println("❌ No available drivers at the moment!")
// 	}
// }

// func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
// 	const R = 6371 // Earth radius in KM

// 	// Convert degrees to radians
// 	dLat := (lat2 - lat1) * (math.Pi / 180)
// 	dLon := (lon2 - lon1) * (math.Pi / 180)

// 	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
// 		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
// 			math.Sin(dLon/2)*math.Sin(dLon/2)

// 	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
// 	distance := R * c

// 	return distance
// }
