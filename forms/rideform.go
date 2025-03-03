package forms

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DriverConnections = make(map[string]*websocket.Conn)
var DriverMutex = sync.Mutex{}
var RideAcceptStatus = make(map[string]bool)

type RideForm struct {
	DriverId         primitive.ObjectID `json:"driver_id" bson:"driver_id" binding:"-"`
	Amount           string             `json:"amount" bson:"amount" binding:"-"`
	StartingPoint    string             `json:"starting_point" bson:"starting_point" binding:"-"`
	DestinationPoint string             `json:"destination_point" bson:"destination_point" binding:"-"`
	StartedTime      *time.Time         `json:"started_time" bson:"started_time" binding:"-"`
	EndTime          *time.Time         `json:"end_time" bson:"end_time" binding:"-"`
}

type DriverLocation struct {
	ID        string  `json:"driver_id" binding:"-"`
	Latitude  float64 `json:"lat" binding:"-"`
	Longitude float64 `json:"lon" binding:"-"`
}

type CustomerLocation struct {
	CustomerId primitive.ObjectID `json:"driver_id" binding:"-"`
	Pickup     PickupLocation     `json:"pickup" binding:"-"`
	Drop       DropLocation       `json:"drop" binding:"-"`
	Radius     float64            `json:"radius"`
}

type PickupLocation struct {
	Latitude  float64 `json:"lat" binding:"-"`
	Longitude float64 `json:"lon" binding:"-"`
}

type DropLocation struct {
	Latitude  float64 `json:"lat" binding:"-"`
	Longitude float64 `json:"lon" binding:"-"`
}

type RideAccept struct {
	RideID   string `json:"ride_id"`
	DriverID string `json:"driver_id"`
	Accepted bool   `json:"accepted"`
}
