package forms

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AccountSid        = "ACbcef2267c53f7781e4259253d31adf9f"
	AuthToken         = "bb0982c09f76a5360e3dc1bd8be62d99"
	FromNumber        = "+1 567 339 2523"
	DB                = "driver"
	DriverCollection  = "driver_info"
	RideCollection    = "ride_info"
	VehicleCollection = "vehicle_info"
)

type UserRegistrationForm struct {
	ProfileImage string `json:"profile_image"` // Base64 string from frontend
	Name         string `json:"name"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Gender       string `json:"gender"`
	KYCData      KYC    `json:"kyc"`
	Password     string `json:"password"`
}

type UserModel struct {
	ProfileImage []byte `json:"profile_image" bson:"profile_image" binding:"-"`
	Name         string `json:"name" bson:"name" binding:"-"`
	Email        string `json:"email" bson:"email" binding:"-"`
	PhoneNumber  string `json:"phone_number" bson:"phone_number" binding:"-"`
	Gender       string `json:"gender" bson:"gender" binding:"-"`
	KYCData      KYC    `json:"kyc" bson:"kyc" binding:"-"`
	Password     string `json:"password" bson:"password" binding:"-"`
}

type LoginUserModel struct {
	PhoneNumber string     `json:"phone_number" bson:"phone_number" binding:"-"`
	Otp         string     `json:"otp" bson:"otp" binding:"-"`
	Time        *time.Time `json:"time" bson:"time" binding:"-"`
}

type KYC struct {
	Aadhaar  string `json:"aadhaar" bson:"aadhaar" binding:"-"`
	DLNumber string `json:"dl_number" bson:"dl_number" binding:"-"`
}

type Vechicle struct {
	DriverId     primitive.ObjectID `json:"driver_id" bson:"driver_id" binding:"-"`
	Registration string             `json:"registration" bson:"registration" binding:"-"`
	Model        string             `json:"model" bson:"model" binding:"-"`
	VehicleImage []byte             `json:"vehicle_image" bson:"vehicle_image" binding:"-"`
}

type LoginModel struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}
