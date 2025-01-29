package models

import (
	"context"
	"nammuru-driver-backend/forms"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type VehicleModel struct{}

func (u *VehicleModel) AddVehicleData(data forms.Vechicle, phoneNumber string) (done string, err error) {
	ctx := context.TODO()
	vehicleCollection := mdbConn.Use(forms.DB, forms.VehicleCollection)
	driverCollection := mdbConn.Use(forms.DB, forms.DriverCollection)

	var profile struct {
		Id primitive.ObjectID `json:"_id" bson:"_id" binding:"-"`
	}

	// Find the driver by phone number
	err = driverCollection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&profile)
	if err != nil {
		return "", err
	}

	data.DriverId = profile.Id

	update := bson.M{"$set": data}

	_, err = vehicleCollection.UpdateOne(ctx, bson.M{"driver_id": profile.Id}, update)
	if err != nil {
		return "", err
	}

	return "Vehicle data added successfully", nil
}

func (u *VehicleModel) AddVehicleImage(data forms.Vechicle) (updateId interface{}, err error) {
	ctx := context.TODO()
	collection := mdbConn.Use(forms.DB, forms.VehicleCollection)

	filter := bson.M{
		"driver_id": data.DriverId,
	}

	update := bson.M{
		"$set": bson.M{
			"vehicle_image": data.VehicleImage,
		},
	}

	updateId, err = collection.UpdateOne(ctx, filter, update)

	if err != nil {
		return updateId, err
	}
	return updateId, err
}
