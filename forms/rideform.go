package forms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideForm struct {
	DriverId         primitive.ObjectID `json:"driver_id" bson:"driver_id"`
	Amount           string             `json:"amount" bson:"amount"`
	StartingPoint    string             `json:"starting_point" bson:"starting_point"`
	DestinationPoint string             `json:"destination_point" bson:"destination_point"`
	StartedTime      *time.Time         `json:"started_time" bson:"started_time"`
	EndTime          *time.Time         `json:"end_time" bson:"end_time"`
}
